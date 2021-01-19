let leaderboard;
let counter = 0;
let page = 1;
const maxEntries = 50;

if (self.fetch) {
    const request = async () => {
        const response = await fetch("/api" + window.location.pathname, {method: 'GET'});
        leaderboard = await response.json();
        renderLeaderboard(1);
    }
    request().then();
} else {
    let request = new XMLHttpRequest();
    request.open("GET", "/api" + window.location.pathname, false);
    request.send();
    leaderboard = JSON.parse(request.responseText);
    renderLeaderboard(1);
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

let counterText = document.createTextNode('#' + (counter) );
span5.appendChild(counterText);
li.appendChild(span5);

return li;


}


function renderLeaderboard(pageNumber) {
                
    let leaderboardlist = document.getElementById("leaderboard-list");
        leaderboardlist.appendChild(document.createTextNode("Use the left/right arrow keys to change page"));




    for (let i = ((pageNumber-1)*(maxEntries)); i < ((pageNumber)*(maxEntries)); i++) {
        
        counter++;
        if(i < leaderboard.length){
        leaderboardlist.appendChild(makeEntry(i));
}
    }
page = pageNumber;

}

document.body.addEventListener("keydown", function (event) {
    

if (event.key === "ArrowRight") {
        
        if(page > 0 && ((page)*(maxEntries)) < leaderboard.length){
        
        clearBox('leaderboard-list');
        renderLeaderboard(page + 1)
    };
    }else if (event.key === "ArrowLeft") {
        
        if(page > 1){
            clearBox('leaderboard-list');
                counter -= 2*maxEntries;
                
            renderLeaderboard(page-1);
}
    }
}, false);



function clearBox(elementID) { 

    var div = document.getElementById(elementID); 



    while(div.firstChild) { 

        div.removeChild(div.firstChild); 

    } 

}
