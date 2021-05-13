function fetchJSON(endpoint) {
    if (self.fetch) {
        return fetch(endpoint, {method: 'GET'})
            .then((resp) => resp.json());
    } else {
        return new Promise(function () {
            let request = new XMLHttpRequest();
            request.open("GET", endpoint, false);
            request.send();
            return JSON.parse(request.responseText);
        })
    }
}

function postJSON(endpoint, data) {
    if (self.fetch) {
        return fetch(endpoint, {
            method: 'POST',
            mode: 'same-origin',
            cache: 'no-cache',
            credentials: 'same-origin',
            redirect: 'error',
            referrerPolicy: 'no-referrer',
            body: JSON.stringify(data),
        })
            .then((resp) => resp.json());
    }
}

function makeOptions() {
    let options = document.createElement("ul");
    options.classList.add("pure-menu-children");
    return options;
}

function makeMenuEntry(text, link) {
    let entry = document.createElement("li");
    entry.classList.add("pure-menu-item");
    let entryLink = document.createElement("a");
    entryLink.classList.add("pure-menu-link");
    entryLink.href = link;
    entryLink.innerText = text;
    entryLink.setAttribute("style", "color: black;");
    entry.appendChild(entryLink);
    return entry;
}
