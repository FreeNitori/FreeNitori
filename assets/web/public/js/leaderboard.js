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

function makeEntry(key){
let entry = leaderboard[key];

let li = document.createElement("LI");
li.setAttribute("class","mdl-list__item mdl-list__item--three-line");
let span1 = document.createElement("Span");
span1.setAttribute("class","mdl-list__item-primary-content");
li.appendChild(span1);

let img1 = document.createElement("img");
img1.setAttribute("src",entry["User"]["AvatarURL"]);
img1.setAttribute("alt","profile_image");
img1.setAttribute("class","mdl-list__item-avatar");
img1.setAttribute("draggable","false");

span1.appendChild(img1);

let span2 = document.createElement("Span");
let span2text = document.createTextNode(entry["User"]["Name"]);
span2.appendChild(span2text);

span1.appendChild(span2);

let span3 = document.createElement("Span");
let span23 = span2.appendChild(span3);
span23.appendChild(document.createTextNode('#' + entry["User"]["Discriminator"]));

let span4 = document.createElement("Span");
span4.setAttribute("class","mdl-list__item-text-body");
span4.appendChild(document.createTextNode("Level: "));
span1.appendChild(span4);


span4.appendChild(document.createTextNode(entry["Level"]));
span4.appendChild(document.createElement("br"));
span4.appendChild(document.createTextNode("Experience: "));


span4.appendChild(document.createTextNode(entry["Experience"]));

let span5 = document.createElement("Span");
span5.setAttribute("class","mdl-list__item-secondary-content");

let counter = document.createTextNode('#' + (key + 1) );
span5.appendChild(counter);
li.appendChild(span5);

return li;


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
    let leaderboardlist = document.getElementById("leaderboard-list");
    for (let i = 0; i < leaderboard.length; i++) {
        leaderboardlist.appendChild(makeEntry(i));
    }
}
