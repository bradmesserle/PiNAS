function startServerMonitor() {

    const monitor = new PollMonitor("/health", {
        mode: "short",                       // or "short" with pollInterval: 5000
        onEvent:   (evt) => console.log(evt.type, evt.data),
        onOnline:  () => console.log("✅ Connected"),
        onOffline: (reason) => console.warn("❌ Disconnected:", reason),
        onRetry:   (n, delay, reason) => console.log(`Retry #${n} in ${delay} ms (${reason})`),
    });
    monitor.start();
}


