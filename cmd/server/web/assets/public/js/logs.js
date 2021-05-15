let logs = document.getElementById("logs");
let ws;
if (window.location.protocol === "https:") {
    ws = new WebSocket("wss://" + window.location.host + "/api/nitori/logs");
} else {
    ws = new WebSocket("ws://" + window.location.host + "/api/nitori/logs");
}
ws.onmessage = function (entry) {
    logs.value += entry.data;
    if (!entry.data.endsWith("\n")) {
        logs.value += "\n";
    }
};
