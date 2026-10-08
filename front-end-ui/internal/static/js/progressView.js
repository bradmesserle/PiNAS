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

function updateText(newValue) {
    let targetDiv = document.getElementById('consoleStream');
    let breakElement = document.createElement('div');
    breakElement.classList.add('break');
    let newSpan = document.createElement('span');
    newSpan.textContent = newValue;
    targetDiv.append(newSpan);
    targetDiv.append(breakElement);
    targetDiv.scrollTop = targetDiv.scrollHeight;

    targetDiv = null;
    breakElement = null;
    newSpan = null;
    newValue = null;
}