/**
 * SSEMonitor — keeps a Server-Sent Events (EventSource) connection alive.
 *
 * The browser's EventSource retries on its own for simple network blips, but it
 * gives up for good (readyState = CLOSED) on things like a 500, a 404, or a wrong
 * Content-Type, and it can't tell when a connection is "open" but silently dead.
 * This class handles all of that:
 *   - reconnects with exponential backoff + jitter whenever the stream closes or errors
 *   - treats the stream as dead if nothing arrives within `heartbeatTimeout`
 *   - optionally pings a health URL first so it doesn't hammer a server that's down
 *   - resumes from the last event ID (sent as a query param, since a new
 *     EventSource can't set the Last-Event-ID header itself)
 *   - reconnects immediately when the browser comes back online or the tab regains focus
 *
 * Browser: works as-is.
 * Node: install the `eventsource` package and pass it in as `EventSourceImpl`.
 *
 * Usage:
 *   const sse = new SSEMonitor("https://your-server.com/events", {
 *     events: {
 *       message: (e) => console.log("message:", e.data),
 *       update:  (e) => console.log("update:", JSON.parse(e.data)),
 *     },
 *     onOpen:    ()         => console.log("✅ Connected"),
 *     onOffline: (reason)   => console.warn("❌ Disconnected:", reason),
 *     onRetry:   (n, delay) => console.log(`Retry #${n} in ${delay} ms`),
 *   });
 *   sse.start();
 *   // later: sse.stop();
 */
class SSEMonitor {
    constructor(url, options = {}) {
        this.url = url;
        this.events = options.events ?? {};              // { eventName: handler(event) }
        this.withCredentials = options.withCredentials ?? false;

        this.healthUrl        = options.healthUrl ?? null; // optional: ping this before reconnecting
        this.healthTimeout    = options.healthTimeout    ?? 5_000;
        this.heartbeatTimeout = options.heartbeatTimeout ?? 45_000; // no data for this long = dead
        this.minRetryDelay    = options.minRetryDelay    ?? 1_000;
        this.maxRetryDelay    = options.maxRetryDelay    ?? 60_000;
        this.maxRetries       = options.maxRetries       ?? Infinity;
        this.lastEventIdParam = options.lastEventIdParam ?? "lastEventId";

        this.onOpen    = options.onOpen    ?? (() => {});
        this.onOffline = options.onOffline ?? (() => {});
        this.onRetry   = options.onRetry   ?? (() => {});
        this.onGiveUp  = options.onGiveUp  ?? (() => {});

        this.EventSourceImpl =
            options.EventSourceImpl ?? (typeof EventSource !== "undefined" ? EventSource : null);
        if (!this.EventSourceImpl) {
            throw new Error("No EventSource available. In Node, pass `EventSourceImpl` (npm i eventsource).");
        }

        this.source = null;
        this.lastEventId = null;
        this.isOnline = false;
        this.retries = 0;
        this.running = false;
        this.retryTimer = null;
        this.heartbeatTimer = null;

        this._onBrowserOnline = () => this._reconnectNow("network back online");
        this._onVisible = () => {
            if (document.visibilityState === "visible" && !this.isOnline) {
                this._reconnectNow("tab visible again");
            }
        };
    }

    start() {
        if (this.running) return;
        this.running = true;
        if (typeof window !== "undefined") {
            window.addEventListener("online", this._onBrowserOnline);
            document.addEventListener("visibilitychange", this._onVisible);
        }
        this._connect();
    }

    stop() {
        this.running = false;
        clearTimeout(this.retryTimer);
        this._teardown();
        if (typeof window !== "undefined") {
            window.removeEventListener("online", this._onBrowserOnline);
            document.removeEventListener("visibilitychange", this._onVisible);
        }
    }

    /* ---------------- connection ---------------- */

    async _connect() {
        if (!this.running) return;
        this._teardown();

        // Optional: confirm the server is up before opening a stream.
        if (this.healthUrl) {
            try {
                await this._ping();
            } catch (err) {
                return this._handleDown(`health check failed: ${err.message}`);
            }
            if (!this.running) return;
        }

        const es = new this.EventSourceImpl(this._buildUrl(), {
            withCredentials: this.withCredentials,
        });
        this.source = es;

        es.onopen = () => {
            if (es !== this.source) return;
            this.isOnline = true;
            this.retries = 0;
            this._resetHeartbeat();
            this.onOpen();
        };

        es.onerror = () => {
            if (es !== this.source) return;
            // Take over reconnection ourselves so backoff and resume logic always apply.
            const reason = es.readyState === 2 ? "connection closed by server" : "connection error";
            this._handleDown(reason);
        };

        // Default unnamed events.
        es.onmessage = (e) => {
            this._track(e);
            this.events.message?.(e);
        };

        // Named events (event: update, event: ping, ...).
        for (const [name, handler] of Object.entries(this.events)) {
            if (name === "message") continue;
            es.addEventListener(name, (e) => {
                this._track(e);
                handler(e);
            });
        }

        // Any event named "ping" or "heartbeat" also keeps the connection alive,
        // even if you didn't register a handler for it.
        for (const name of ["ping", "heartbeat"]) {
            if (!this.events[name]) es.addEventListener(name, (e) => this._track(e));
        }
    }

    _track(e) {
        if (e.lastEventId) this.lastEventId = e.lastEventId;
        this._resetHeartbeat();
    }

    _buildUrl() {
        if (!this.lastEventId) return this.url;
        const base = typeof location !== "undefined" ? location.href : undefined;
        const u = new URL(this.url, base);
        u.searchParams.set(this.lastEventIdParam, this.lastEventId);
        return u.toString();
    }

    _teardown() {
        clearTimeout(this.heartbeatTimer);
        if (this.source) {
            this.source.onopen = this.source.onerror = this.source.onmessage = null;
            this.source.close();
            this.source = null;
        }
    }

    /* ---------------- liveness ---------------- */

    _resetHeartbeat() {
        clearTimeout(this.heartbeatTimer);
        if (!this.heartbeatTimeout) return;
        this.heartbeatTimer = setTimeout(
            () => this._handleDown(`no data for ${this.heartbeatTimeout} ms`),
            this.heartbeatTimeout
        );
    }

    async _ping() {
        const controller = new AbortController();
        const t = setTimeout(() => controller.abort(), this.healthTimeout);
        try {
            const res = await fetch(this.healthUrl, { cache: "no-store", signal: controller.signal });
            if (!res.ok) throw new Error(`HTTP ${res.status}`);
        } catch (err) {
            if (err.name === "AbortError") throw new Error(`timed out after ${this.healthTimeout} ms`);
            throw err;
        } finally {
            clearTimeout(t);
        }
    }

    /* ---------------- reconnect ---------------- */

    _handleDown(reason) {
        this._teardown();
        if (!this.running) return;

        if (this.isOnline) {
            this.isOnline = false;
            this.onOffline(reason);
        }

        if (this.retries >= this.maxRetries) {
            this.onGiveUp(this.retries);
            this.stop();
            return;
        }

        // Exponential backoff with jitter: ~1s, 2s, 4s, 8s ... capped at maxRetryDelay.
        const base = Math.min(this.minRetryDelay * 2 ** this.retries, this.maxRetryDelay);
        const delay = Math.round(base / 2 + Math.random() * (base / 2));
        this.retries++;
        this.onRetry(this.retries, delay, reason);

        clearTimeout(this.retryTimer);
        this.retryTimer = setTimeout(() => this._connect(), delay);
    }

    _reconnectNow(reason) {
        if (!this.running || this.isOnline) return;
        clearTimeout(this.retryTimer);
        this.retries = 0;
        console.info(`Reconnecting now (${reason})`);
        this._connect();
    }
}


/* ---------------- Example ---------------- */
// Browser:
// const sse = new SSEMonitor("/events", {
//   healthUrl: "/health",              // optional
//   heartbeatTimeout: 45000,           // server should send something at least every ~30s
//   events: {
//     message: (e) => console.log("message:", e.data),
//     update:  (e) => console.log("update:", JSON.parse(e.data)),
//   },
//   onOpen:    () => console.log("✅ Connected"),
//   onOffline: (reason) => console.warn("❌ Disconnected:", reason),
//   onRetry:   (n, delay, reason) => console.log(`Retry #${n} in ${delay} ms (${reason})`),
// });
// sse.start();
//
// Node:
// const { EventSource } = require("eventsource");
// const sse = new SSEMonitor("https://your-server.com/events", { EventSourceImpl: EventSource, ... });
//
// Server side, send a keep-alive so the heartbeat check has something to see, e.g. every 30s:
//   res.write("event: ping\ndata: {}\n\n");