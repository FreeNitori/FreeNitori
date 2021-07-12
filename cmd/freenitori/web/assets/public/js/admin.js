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

    // Mem
    document.getElementById("MemAllocated").innerText = data["Mem"]["Allocated"];
    document.getElementById("MemTotal").innerText = data["Mem"]["Total"];
    document.getElementById("MemSys").innerText = data["Mem"]["Sys"];
    document.getElementById("MemLookups").innerText = data["Mem"]["Lookups"];
    document.getElementById("MemMallocs").innerText = data["Mem"]["Mallocs"];
    document.getElementById("MemFrees").innerText = data["Mem"]["Frees"];

    // Heap
    document.getElementById("HeapAlloc").innerText = data["Heap"]["Alloc"];
    document.getElementById("HeapSys").innerText = data["Heap"]["Sys"];
    document.getElementById("HeapIdle").innerText = data["Heap"]["Idle"];
    document.getElementById("HeapInuse").innerText = data["Heap"]["Inuse"];
    document.getElementById("HeapReleased").innerText = data["Heap"]["Released"];
    document.getElementById("HeapObjects").innerText = data["Heap"]["Objects"];

    // GC
    document.getElementById("NextGC").innerText = data["GC"]["NextGC"];
    document.getElementById("LastGC").innerText = data["GC"]["LastGC"];
    document.getElementById("PauseTotalNs").innerText = data["GC"]["PauseTotalNs"];
    document.getElementById("PauseNs").innerText = data["GC"]["PauseNs"];
    document.getElementById("NumGC").innerText = data["GC"]["NumGC"];

    // Misc
    document.getElementById("StackInuse").innerText = data["Misc"]["StackInuse"];
    document.getElementById("StackSys").innerText = data["Misc"]["StackSys"];
    document.getElementById("MSpanInuse").innerText = data["Misc"]["MSpanInuse"];
    document.getElementById("MSpanSys").innerText = data["Misc"]["MSpanSys"];
    document.getElementById("MCacheInuse").innerText = data["Misc"]["MCacheInuse"];
    document.getElementById("MCacheSys").innerText = data["Misc"]["MCacheSys"];
    document.getElementById("GCSys").innerText = data["Misc"]["GCSys"];
    document.getElementById("BuckHashSys").innerText = data["Misc"]["BuckHashSys"];
    document.getElementById("OtherSys").innerText = data["Misc"]["OtherSys"];
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
