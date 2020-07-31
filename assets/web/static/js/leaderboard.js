let request = new XMLHttpRequest();
request.open("GET", "/api" + window.location.pathname, false);
request.send();
let response = JSON.parse(request.responseText)

function makeEntry(key) {
    let entry = response[key]
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
        `
}

for (let i = 0; i < response.length; i++) {
    document.getElementById("leaderboard-body").innerHTML += makeEntry(i)
}
// $(document).ready(function () {
//     $('#leaderboard-body').endlessScroll({
//         fireOnce: false,
//         callback: function () {}
//     })
// })