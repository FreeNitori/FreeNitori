document.getElementById("snowflake").addEventListener("submit", function (event) {
    lookup(document.getElementById('snowflake').value);
    event.preventDefault();
}, false);

function lookup(snowflake) {
    let userInfo = fetchJSON("http://localhost:7777/api/user/" + snowflake)
    userInfo.then(function (data) {
        document.getElementById("result").innerHTML = data.text;
    })
}