let leaderboard;
let guildname;
let guildicon;
let counter = 0;
let page = 1;
const maxEntries = 50;

if (self.fetch) {
    const request = async () => {
        const response = await fetch("/api" + window.location.pathname, {method: 'GET'});
        leaderboard = await response.json();
        HeadAppend();
        layoutTitle();
        renderLeaderboard(1);
    }
    request().then();
} else {
    let request = new XMLHttpRequest();
    request.open("GET", "/api" + window.location.pathname, false);
    request.send();
    leaderboard = JSON.parse(request.responseText);
    HeadAppend();
    renderLeaderboard(1);
}


function makeEntry(key) {
    let entry = leaderboard["Leaderboard"][key];

    let li = document.createElement("LI");
    li.setAttribute("class", "mdl-list__item mdl-list__item--three-line");
    let span1 = document.createElement("Span");
    span1.setAttribute("class", "mdl-list__item-primary-content");
    li.appendChild(span1);

    let img1 = document.createElement("img");
    img1.setAttribute("src", entry["User"]["AvatarURL"]);
    img1.setAttribute("alt", "profile_image");
    img1.setAttribute("class", "mdl-list__item-avatar");
    img1.setAttribute("draggable", "false");

    span1.appendChild(img1);

    let span2 = document.createElement("Span");
    let span2text = document.createTextNode(entry["User"]["Name"]);
    span2.appendChild(span2text);

    span1.appendChild(span2);

    let span3 = document.createElement("Span");
    let span23 = span2.appendChild(span3);
    span23.appendChild(document.createTextNode('#' + entry["User"]["Discriminator"]));

    let span4 = document.createElement("Span");
    span4.setAttribute("class", "mdl-list__item-text-body");
    span4.appendChild(document.createTextNode("Level: "));
    span1.appendChild(span4);


    span4.appendChild(document.createTextNode(entry["Level"]));
    span4.appendChild(document.createElement("br"));
    span4.appendChild(document.createTextNode("Experience: "));


    span4.appendChild(document.createTextNode(entry["Experience"]));

    let span5 = document.createElement("Span");
    span5.setAttribute("class", "mdl-list__item-secondary-content");

    let counterText = document.createTextNode('#' + (counter));
    span5.appendChild(counterText);
    li.appendChild(span5);

    return li;


}


function renderLeaderboard(pageNumber) {

    let leaderboardlist = document.getElementById("leaderboard-list");
    let fragment = document.createDocumentFragment();


    for (let i = ((pageNumber - 1) * (maxEntries)); i < ((pageNumber) * (maxEntries)); i++) {

        counter++;
        if (i < leaderboard.Leaderboard.length) {
            fragment.appendChild(makeEntry(i));

        }
    }
    leaderboardlist.appendChild(fragment);
    page = pageNumber;

}

document.body.addEventListener("keydown", function (event) {


    if (event.key === "ArrowRight") {

        if (page > 0 && ((page) * (maxEntries)) < leaderboard.length) {

            clearBox('leaderboard-list');
            renderLeaderboard(page + 1)
        }
        ;
    } else if (event.key === "ArrowLeft") {

        if (page > 1) {
            clearBox('leaderboard-list');
            counter -= 2 * maxEntries;

            renderLeaderboard(page - 1);
        }
    }
}, false);


function clearBox(elementID) {

    var div = document.getElementById(elementID);


    while (div.firstChild) {

        div.removeChild(div.firstChild);

    }

}


function HeadAppend() {
    guildname = leaderboard["GuildInfo"]["Name"];
    guildicon = leaderboard["GuildInfo"]["IconURL"];

    let title = document.createElement("title");
    titleText = document.createTextNode(guildname + ' Leaderboard');
    title.appendChild(titleText);
    (document.head).appendChild(title);

    let HTTPHeader = document.createElement("meta");
    HTTPHeader.setAttribute("http-equiv", "Content-Type");
    HTTPHeader.setAttribute("content", "text/html; charset=UTF-8");
    (document.head).appendChild(HTTPHeader);

    let viewportHeader = document.createElement("meta");
    viewportHeader.setAttribute("name", "viewport");
    viewportHeader.setAttribute("content", "width=device-width, initial-scale=1, maximum-scale=1.0, user-scalable=no");
    (document.head).appendChild(viewportHeader);

    let RobotHeader = document.createElement("meta");
    RobotHeader.setAttribute("name", "robot");
    RobotHeader.setAttribute("content", "noindex, nofollow");
    (document.head).appendChild(RobotHeader);

    let guildIconLink = document.createElement("link");
    guildIconLink.setAttribute("rel", "shortcut icon");
    guildIconLink.setAttribute("href", guildicon);
    (document.head).appendChild(guildIconLink);

    let guildIconLink2 = document.createElement("link");
    guildIconLink2.setAttribute("rel", "icon");
    guildIconLink2.setAttribute("href", guildicon);
    (document.head).appendChild(guildIconLink2);

    let meta1 = document.createElement("meta");
    meta1.setAttribute("property", "og:title");
    meta1.setAttribute("content", guildname + ' Leaderboard');
    (document.head).appendChild(meta1);

    let meta2 = document.createElement("meta");
    meta2.setAttribute("property", "og:description");
    meta2.setAttribute("content", 'Experience Leaderboard of ' + guildname);
    (document.head).appendChild(meta2);

    let meta3 = document.createElement("meta");
    meta3.setAttribute("property", "og:type");
    meta3.setAttribute("content", "website");
    (document.head).appendChild(meta3);

    let meta4 = document.createElement("meta");
    meta4.setAttribute("property", "og:url");
    meta4.setAttribute("content", "/");
    (document.head).appendChild(meta4);

    let meta5 = document.createElement("meta");
    meta5.setAttribute("property", "og:image");
    meta5.setAttribute("content", guildicon);
    (document.head).appendChild(meta5);

    let meta6 = document.createElement("meta");
    meta6.setAttribute("property", "og:site_name");
    meta6.setAttribute("content", guildname + ' Leaderboard');
    (document.head).appendChild(meta6);

    let meta7 = document.createElement("meta");
    meta7.setAttribute("itemprop", "name");
    meta7.setAttribute("content", guildname + ' Leaderboard');
    (document.head).appendChild(meta7);

    let meta8 = document.createElement("meta");
    meta8.setAttribute("itemdrop", "description");
    meta8.setAttribute("content", "Experience leaderboard of " + guildname);
    (document.head).appendChild(meta8);

    let meta9 = document.createElement("meta");
    meta9.setAttribute("name", "twitter:card");
    meta9.setAttribute("content", "summary_large_image");
    (document.head).appendChild(meta9);

    let meta10 = document.createElement("meta");
    meta10.setAttribute("name", "twitter:image");
    meta10.setAttribute("content", guildicon);
    (document.head).appendChild(meta10);

    let meta11 = document.createElement("meta");
    meta10.setAttribute("name", "twitter:title");
    meta10.setAttribute("content", guildname + ' Leaderboard');
    (document.head).appendChild(meta11);

    meta12 = document.createElement("meta");
    meta12.setAttribute("name", "twitter:description");
    meta12.setAttribute("content", 'Experience Leaderboard of ' + guildname);
    (document.head).appendChild(meta12);

    let guildiconID = document.getElementsByClassName('guild-icon')[0];
    guildiconID.style.background = 'url(' + guildicon + ')  50% no-repeat';
};

function layoutTitle() {
    let title = document.getElementsByClassName("mdl-layout__title")[0];

    title.appendChild(document.createTextNode('Leaderboard of ' + guildname));
}