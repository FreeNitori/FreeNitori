let leaderboard;

if (self.fetch) {
    const request = async () => {
        const response = await fetch("/api" + window.location.pathname, {method: 'GET'});
        leaderboard = await response.json();
        renderLeaderboard();
    }
    request().then();
} else {
    let request = new XMLHttpRequest();
    request.open("GET", "/api" + window.location.pathname, false);
    request.send();
    leaderboard = JSON.parse(request.responseText);
    renderLeaderboard();
}

function makeEntry(key) {
    let entry = leaderboard[key];
    return `
                        <li class="mdl-list__item mdl-list__item--three-line">
                          <span class="mdl-list__item-primary-content">
                          <img src="` + entry["User"]["AvatarURL"] + `"
                               alt="profile_image" class="mdl-list__item-avatar" draggable="false">
                          <span>` + entry["User"]["Name"] + `<span style="color: #808080; ">#` + entry["User"]["Discriminator"] + `</span></span>
                          <span class="mdl-list__item-text-body">
                          Level: <a>` + entry["Level"] + `</a>
                          <br>
                          Experience: <a>` + entry["Experience"] + `</a>
                          </span>
                          </span>
                          <span class="mdl-list__item-secondary-content">
                             <a class="mdl-list__item-secondary-action">#` + (key + 1) + `</a>
                          </span>
                        </li>
        `;
}

function makePage(index) {
    let page = "";
    let finalKey = 9;
    if (leaderboard.length > index * 10) {
        if (leaderboard.length > (index + 1) * 10) {
            finalKey = (index + 1) * 10;
        } else {
            finalKey = leaderboard.length;
        }
        for (let i = index * 10; i < finalKey; i++) {
            page += makeEntry(i);
        }
        return page;
    } else {
        return null;
    }
}

function renderLeaderboard() {
    for (let i = 0; i < leaderboard.length; i++) {
        document.getElementById("leaderboard-list").innerHTML += makeEntry(i);
    }
}