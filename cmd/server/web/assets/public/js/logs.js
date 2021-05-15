let logs = document.getElementById("logs");
let ws;
openSocket();

function openSocket() {
    if (window.location.protocol === "https:") {
        ws = new WebSocket("wss://" + window.location.host + "/api/nitori/logs");
    } else {
        ws = new WebSocket("ws://" + window.location.host + "/api/nitori/logs");
    }

    ws.onopen = function (_) {
        console.log("Socked opened.");
    };

    ws.onmessage = function (event) {
        logs.value += event.data;
        if (!event.data.endsWith("\n")) {
            logs.value += "\n";
        }
    };

    ws.onclose = function (_) {
        console.log("Socket close.");
        setTimeout(function () {
            openSocket();
        }, 1000);
    };

    ws.onerror = function (_) {
        console.error("Socket error.");
        ws.close();
    }
}