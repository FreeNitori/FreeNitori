document.getElementById("snowflake").addEventListener("keydown", function (event) {
    if (event.key === "Enter") {
	clearBox('result');
        lookup(document.getElementById('snowflake').value);
    }
}, false);

function lookup(snowflake) {
    let userInfo = fetchJSON("http://localhost:7777/api/user/" + snowflake)
    userInfo.then(function (data) {

        let userData = document.getElementById("result");
        let ul = document.createElement("ul");

        for (i in data) {
            let li = document.createElement("LI");
if(data.AvatarURL == data[i]){
let userAvatar = document.createElement("IMG");
userAvatar.setAttribute("src",data.AvatarURL);
userAvatar.setAttribute("alt","data.AvatarURL");

li.appendChild(userAvatar);
}else{
            let userText = document.createTextNode(i + ': ' + data[i]);
            li.appendChild(userText);
}
            ul.appendChild(li);
        }

        userData.appendChild(ul);

    })
}


function clearBox(elementID) { 

            var div = document.getElementById(elementID); 



            while(div.firstChild) { 

                div.removeChild(div.firstChild); 

            } 

        }