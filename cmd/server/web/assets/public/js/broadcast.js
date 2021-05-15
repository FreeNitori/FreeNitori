// Populate from /api/nitori/broadcast
fetchJSON("/api/nitori/broadcast").then(function (data) {
    document.getElementById("broadcast").innerText = data["content"];
});

function sendBroadcast(silent) {
    postJSON("/api/nitori/broadcast", {
        "alert": !silent,
        "content": document.getElementById("broadcast").value,
    }).then(function (data) {
        if (data["state"] === "ok") {
            if (silent) {
                alert("Successfully broadcast message silently.");
            } else {
                alert("Successfully broadcast message.");
            }
        } else {
            alert("Error while broadcasting message: " + data["error"]);
        }
    });
}

function clearBroadcastBuffer() {
    document.getElementById("broadcast").innerText = "";
    sendBroadcast(true);
}
