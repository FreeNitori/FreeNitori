// I don't know JavaScript, so this is probably a very terrible script, if you are interested in contributing, hit me up on Discord or something :)
let info = fetchJSON("/api/info")
let stats = fetchJSON("/api/stats")

info.then(function (data) {
    document.getElementById("inviteURL").href = data["invite_url"];
    document.getElementById("programVersion").textContent = data["nitori_version"];
})
stats.then(function (data) {
    document.getElementById("messageTotal").textContent = data["total_messages"];
    document.getElementById("guildsDeployed").textContent = data["guilds_deployed"];
})

function fetchJSON(endpoint) {
    if (self.fetch) {
        return fetch(endpoint, {method: 'GET'})
            .then((resp) => resp.json());
    } else {
        return new Promise(function () {
            let request = new XMLHttpRequest();
            request.open("GET", endpoint, false);
            request.send();
            return JSON.parse(request.responseText);
        })
    }
}