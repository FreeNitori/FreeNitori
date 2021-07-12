// Populate from /api/info
fetchJSON("/api/info").then(function (data) {
    document.getElementById("inviteURL").href = data["invite_url"];
    document.getElementById("programVersion").textContent = data["nitori_version"];
    document.getElementById("programRevision").textContent = data["nitori_revision"];
});

// Populate from /api/stats
fetchJSON("/api/stats").then(function (data) {
    document.getElementById("messageTotal").textContent = data["total_messages"];
    document.getElementById("guildsDeployed").textContent = data["guilds_deployed"];
});

// Wipe results area
function wipeResults() {
    document.getElementById("lookupResultTitle0").textContent = "";
    document.getElementById("lookupResultContent0").textContent = "";
    document.getElementById("lookupResultTitle1").textContent = "";
    document.getElementById("lookupResultContent1").textContent = "";
    document.getElementById("lookupResultTitle2").textContent = "";
    document.getElementById("lookupResultContent2").textContent = "";
    document.getElementById("lookupResultTitle3").textContent = "";
    document.getElementById("lookupResultContent3").textContent = "";
}

// Lookup user
function lookupUser() {
    wipeResults();
    let snowflake = document.getElementById("userLookupField").value;
    if (parseInt(snowflake) >>> 22 <= 0) {
        document.getElementById("lookupResultTitle0").textContent = "Error";
        document.getElementById("lookupResultContent0").textContent = "Invalid snowflake.";
        return null;
    }
    fetchJSON("/api/user/" + snowflake).then(function (data) {
        if (data["error"] != null) {
            document.getElementById("lookupResultTitle0").textContent = "Error";
            document.getElementById("lookupResultContent0").textContent = data["error"];
            return;
        }
        document.getElementById("lookupResultTitle0").textContent = "Username";
        document.getElementById("lookupResultContent0").innerHTML = data["Name"] + `<span style="color: gray">#</span>` + data["Discriminator"];
        document.getElementById("lookupResultTitle1").textContent = "Creation Time";
        document.getElementById("lookupResultContent1").innerText = data["CreationTime"];
        document.getElementById("lookupResultTitle2").textContent = "Profile Picture";
        document.getElementById("lookupResultContent2").innerHTML = `<a href="` + data["AvatarURL"] + `" target="_blank">Open</a>`;
        document.getElementById("lookupResultTitle3").textContent = "Bot User";
        document.getElementById("lookupResultContent3").innerText = data["Bot"];
    });
}

// Lookup guild
function lookupGuild() {
    wipeResults();
    let snowflake = document.getElementById("guildLookupField").value;
    if (parseInt(snowflake) >>> 22 <= 0) {
        document.getElementById("lookupResultTitle0").textContent = "Error";
        document.getElementById("lookupResultContent0").textContent = "Invalid snowflake.";
        return null;
    }
    fetchJSON("/api/guild/" + snowflake).then(function (data) {
        if (data["error"] != null) {
            document.getElementById("lookupResultTitle0").textContent = "Error";
            document.getElementById("lookupResultContent0").textContent = data["error"];
            return;
        }
        document.getElementById("lookupResultTitle0").textContent = "Name";
        document.getElementById("lookupResultContent0").innerHTML = data["Name"];
        document.getElementById("lookupResultTitle1").textContent = "Creation Time";
        document.getElementById("lookupResultContent1").innerText = data["CreationTime"];
        document.getElementById("lookupResultTitle2").textContent = "Icon";
        document.getElementById("lookupResultContent2").innerHTML = `<a href="` + data["IconURL"] + `" target="_blank">Open</a>`;
        document.getElementById("lookupResultTitle3").textContent = "Member Count";
        document.getElementById("lookupResultContent3").innerText = data["Members"].length;
    });
}
