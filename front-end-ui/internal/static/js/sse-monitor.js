/**
 * PollMonitor — receives updates from a server and stays connected using
 * plain HTTP requests only (fetch). No Server-Sent Events, no WebSockets.
 *
 * Two modes:
 *   "long"  (default) — long polling: the server holds each request open until it
 *            has new data (or ~30s pass), so updates arrive almost instantly.
 *   "short" — asks the server every `pollInterval` ms whether anything is new.
 *            Simpler server side, slightly delayed updates.
 *
 * If a request fails or times out, the server is treated as offline and the
 * monitor retries with exponential backoff + jitter until it's back. It also
 * retries immediately when the browser comes back online or the tab regains focus.
 *
 * Works in the browser and in Node 18+ (built-in fetch + AbortController).
 *
 * Expected server response (JSON, status 200):
 *   { "cursor": "123", "events": [ { "id": "123", "type": "update", "data": {...} } ] }
 * A 204 No Content means "nothing new" and is fine too.
 * The monitor sends the last cursor back as ?cursor=123 so nothing is missed.
 *
 * Usage:
 *   const monitor = new PollMonitor("https://your-server.com/updates", {
 *     onEvent:   (evt)      => console.log(evt.type, evt.data),
 *     onOnline:  ()         => console.log("✅ Connected"),
 *     onOffline: (reason)   => console.warn("❌ Disconnected:", reason),
 *     onRetry:   (n, delay) => console.log(`Retry #${n} in ${delay} ms`),
 *   });
 *   monitor.start();
 *   // later: monitor.stop();
 */
class PollMonitor {
    constructor(url, options = {}) {
        this.url = url;
        this.mode          = options.mode          ?? "long";  // "long" | "short"
        this.pollInterval  = options.pollInterval  ?? 5_000;   // short mode: ms between polls
        this.requestTimeout =
            options.requestTimeout ?? (this.mode === "long" ? 40_000 : 5_000); // long: > server hold time
        this.minRetryDelay = options.minRetryDelay ?? 1_000;
        this.maxRetryDelay = options.maxRetryDelay ?? 60_000;
        this.maxRetries    = options.maxRetries    ?? Infinity;
        this.cursorParam   = options.cursorParam   ?? "cursor";
        this.fetchOptions  = options.fetchOptions  ?? {};       // extra fetch options, e.g. headers

        this.onEvent   = options.onEvent   ?? (() => {});
        this.onOnline  = options.onOnline  ?? (() => {});
        this.onOffline = options.onOffline ?? (() => {});
        this.onRetry   = options.onRetry   ?? (() => {});
        this.onGiveUp  = options.onGiveUp  ?? (() => {});

        this.cursor = options.cursor ?? null;
        this.isOnline = false;
        this.retries = 0;
        this.running = false;
        this.timer = null;
        this.controller = null;

        this._onBrowserOnline = () => this._retryNow("network back online");
        this._onVisible = () => {
            if (document.visibilityState === "visible") this._retryNow("tab visible again");
        };
    }

    start() {
        if (this.running) return;
        this.running = true;
        if (typeof window !== "undefined") {
            window.addEventListener("online", this._onBrowserOnline);
            document.addEventListener("visibilitychange", this._onVisible);
        }
        this._poll();
    }

    stop() {
        this.running = false;
        clearTimeout(this.timer);
        this.controller?.abort();
        if (typeof window !== "undefined") {
            window.removeEventListener("online", this._onBrowserOnline);
            document.removeEventListener("visibilitychange", this._onVisible);
        }
    }

    /* ---------------- polling loop ---------------- */

    async _poll() {
        if (!this.running) return;
        clearTimeout(this.timer);

        const controller = new AbortController();
        this.controller = controller;
        const timeout = setTimeout(() => controller.abort("timeout"), this.requestTimeout);

        try {
            const res = await fetch(this._buildUrl(), {
                cache: "no-store",
                ...this.fetchOptions,
                signal: controller.signal,
            });
            if (res.status !== 200 && res.status !== 204) throw new Error(`HTTP ${res.status}`);

            const body = res.status === 204 ? null : await res.json();
            if (controller !== this.controller || !this.running) return; // superseded or stopped

            this._markOnline();
            this._handleBody(body);

            // Long polling: ask again right away. Short polling: wait first.
            this._schedule(this.mode === "long" ? 0 : this.pollInterval);
        } catch (err) {
            if (controller !== this.controller || !this.running) return;
            const reason =
                controller.signal.aborted ? `no response after ${this.requestTimeout} ms` : err.message;
            this._handleDown(reason);
        } finally {
            clearTimeout(timeout);
        }
    }

    _handleBody(body) {
        if (!body) return;
        for (const evt of body.events ?? []) {
            try {
                this.onEvent(evt);
            } catch (err) {
                console.error("onEvent handler threw:", err);
            }
        }
        const last = body.events?.at(-1);
        this.cursor = body.cursor ?? last?.id ?? this.cursor;
    }

    _buildUrl() {
        if (this.cursor == null) return this.url;
        const base = typeof location !== "undefined" ? location.href : undefined;
        const u = new URL(this.url, base);
        u.searchParams.set(this.cursorParam, this.cursor);
        return u.toString();
    }

    _schedule(ms) {
        if (!this.running) return;
        clearTimeout(this.timer);
        this.timer = setTimeout(() => this._poll(), ms);
    }

    /* ---------------- online / offline ---------------- */

    _markOnline() {
        this.retries = 0;
        if (!this.isOnline) {
            this.isOnline = true;
            this.onOnline();
        }
    }

    _handleDown(reason) {
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

        this._schedule(delay);
    }

    _retryNow(reason) {
        if (!this.running || this.isOnline) return;
        console.info(`Reconnecting now (${reason})`);
        this.retries = 0;
        this.controller?.abort();
        this.controller = null;
        this._poll();
    }
}

if (typeof module !== "undefined" && module.exports) {
    module.exports = PollMonitor;
}

/* ---------------- Example ---------------- */
// const monitor = new PollMonitor("/updates", {
//   mode: "long",                       // or "short" with pollInterval: 5000
//   onEvent:   (evt) => console.log(evt.type, evt.data),
//   onOnline:  () => console.log("✅ Connected"),
//   onOffline: (reason) => console.warn("❌ Disconnected:", reason),
//   onRetry:   (n, delay, reason) => console.log(`Retry #${n} in ${delay} ms (${reason})`),
// });
// monitor.start();
//
// Minimal Express long-poll endpoint for reference:
//   app.get("/updates", (req, res) => {
//     const since = req.query.cursor;
//     const pending = getEventsAfter(since);           // your own storage
//     if (pending.length) return res.json({ cursor: pending.at(-1).id, events: pending });
//     const done = (events) => res.json({ cursor: events.at(-1).id, events });
//     waiters.add(done);                              // call done(newEvents) when data arrives
//     req.on("close", () => waiters.delete(done));
//     setTimeout(() => { waiters.delete(done); if (!res.headersSent) res.status(204).end(); }, 30000);
//   });