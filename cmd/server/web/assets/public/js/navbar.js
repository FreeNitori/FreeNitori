// Populate from /api/auth
fetchJSON("/api/auth").then(function (data) {
    let button = document.getElementById("oauthButton");
    if (!data["authorized"]) {
        button.innerText = "Login";
        button.href = "/auth/login";
        button.setAttribute("style", "color: white;");
    } else {
        button.href = "#";
        button.onclick = function () {
            return false;
        };
        fetchJSON("/api/auth/user").then(function (data) {
            let menu = document.getElementById("oauthMenu");
            menu.classList.add("pure-menu-has-children");
            menu.classList.add("pure-menu-allow-hover");
            button.innerText = data["user"]["Name"];
            let options = makeOptions();
            if (data["administrator"]) {
                options.appendChild(makeMenuEntry("Administration", "/auth/admin"));
            } else if (data["operator"]) {
                options.appendChild(makeMenuEntry("Administration", "/auth/operator"));
            }
            options.appendChild(makeMenuEntry("Logout", "/auth/logout"));
            menu.appendChild(options);
        });
    }
});
