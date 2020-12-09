document.getElementById("snowflake").addEventListener("keydown", function (event) {
    if (event.key === "Enter") {
        lookup(document.getElementById('snowflake').value);
    }
}, false);

function lookup(snowflake) {
    let userInfo = fetchJSON("http://localhost:7777/api/user/" + snowflake)
    userInfo.then(function (data) {
        document.getElementById("result").innerHTML = data["Name"];
    })
}