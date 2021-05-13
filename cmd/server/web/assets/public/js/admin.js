// Populate from /api/nitori
fetchJSON("/api/nitori").then(function (data) {
    document.getElementById("newUsername").value = data["Name"];
    document.getElementById("newDiscriminator").value = data["Discriminator"];
});

function changeUsername() {}