let nitoriData;

// Populate from /api/nitori
fetchJSON("/api/nitori").then(function (data) {
    document.getElementById("newUsername").value = data["Name"];
    document.getElementById("discriminator").innerText = data["Discriminator"];
    nitoriData = data;
});

// Populate from /api/nitori/stats
fetchJSON("/api/nitori/stats").then(function (data) {
    // Process
    document.getElementById("PID").innerText = data["Process"]["PID"];
    document.getElementById("Uptime").innerText = data["Process"]["Uptime"];
    document.getElementById("NumGoroutine").innerText = data["Process"]["NumGoroutine"];
    document.getElementById("DBSize").innerText = data["Process"]["DBSize"];

    // Platform
    document.getElementById("GoVersion").innerText = data["Platform"]["GoVersion"];
    document.getElementById("GOOS").innerText = data["Platform"]["GOOS"];
    document.getElementById("GOARCH").innerText = data["Platform"]["GOARCH"];
    document.getElementById("GOROOT").innerText = data["Platform"]["GOROOT"];

    // Discord
    document.getElementById("Intents").innerText = data["Discord"]["Intents"];
    document.getElementById("Sharding").innerText = data["Discord"]["Sharding"];
    document.getElementById("Shards").innerText = data["Discord"]["Shards"];
    document.getElementById("Guilds").innerText = data["Discord"]["Guilds"];
});

function changeUsername() {
    nitoriData["Name"] = document.getElementById("newUsername").value;
    postJSON("/api/nitori", nitoriData).then(function (data) {
        if (data["state"] === "ok") {
            alert("Successfully updated username.");
        } else {
            alert("Error while updating username: " + data["error"]);
        }
    });
}

function executeAction(action) {
    let actionPayload = {
        "action": action,
    }
    postJSON("/api/nitori/action", actionPayload).then(function (data) {
        if (data["state"] === "ok") {
            alert("Successfully performed action: " + action);
        } else {
            alert("Error while performing action: " + data["error"]);
        }
    });
}
