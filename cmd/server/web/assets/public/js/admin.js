let nitoriData;

// Populate from /api/nitori
fetchJSON("/api/nitori").then(function (data) {
    document.getElementById("newUsername").value = data["Name"];
    document.getElementById("discriminator").innerText = data["Discriminator"];
    nitoriData = data;
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
