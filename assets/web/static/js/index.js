let stats;

if (self.fetch) {
    const request = async () => {
        const response = await fetch("/api/stats", {method: 'GET'});
        stats = await response.json();
        finalizePage();
    }
    request().then();
} else {
    let request = new XMLHttpRequest();
    request.open("GET", "/api/stats", false);
    request.send();
    leaderboard = JSON.parse(request.responseText);
    finalizePage();
}

function finalizePage() {
    document.getElementById("messageTotal").textContent = stats["total_messages"];
    document.getElementById("guildsDeployed").textContent = stats["guilds_deployed"];
    document.getElementById("programVersion").textContent = stats["nitori_version"];
}

function redirectInvite() {
    if (self.fetch) {
        fetch("/api/invite", {method: 'GET'}).then(
            response => response.json()
        ).then(
            function (response) {
                window.open(response["invite_url"]);
            }
        )
    } else {
        let request = new XMLHttpRequest();
        request.open("GET", "/api/invite", false);
        request.send();
        let response = JSON.parse(request.responseText);
        window.open(response["invite_url"]);
    }
}