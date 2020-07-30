let request = new XMLHttpRequest();
request.open("GET", "/api/stats", false);
request.send();

let stats = JSON.parse(request.responseText)
document.getElementById("messageTotal").textContent = stats["total_messages"]
document.getElementById("guildsDeployed").textContent = stats["guilds_deployed"]
document.getElementById("programVersion").textContent = stats["nitori_version"]

function redirectInvite() {
    let request = new XMLHttpRequest();
    request.open("GET", "/api/invite", false);
    request.send();
    let response = JSON.parse(request.responseText)
    window.open(response["invite_url"])
}