// This script works but feels quite terrible, contribute something if able.
let info = fetchJSON("/api/info")
let stats = fetchJSON("/api/stats")

info.then(function (data) {
    document.getElementById("inviteURL").href = data["invite_url"];
    document.getElementById("programVersion").textContent = data["nitori_version"];
    document.getElementById("programRevision").textContent = data["nitori_revision"];
})
stats.then(function (data) {
    document.getElementById("messageTotal").textContent = data["total_messages"];
    document.getElementById("guildsDeployed").textContent = data["guilds_deployed"];
})