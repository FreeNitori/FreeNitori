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