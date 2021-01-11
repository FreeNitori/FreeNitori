document.getElementById("snowflake").addEventListener("keydown", function (event) {
    if (event.key === "Enter") {
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
            let userText = document.createTextNode(data[i]);
            li.appendChild(userText);
            ul.appendChild(li);
        }

        userData.appendChild(ul);

    })
}
