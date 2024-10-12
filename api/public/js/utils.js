String.prototype.ohMySimpleHash = function () {
    let hash = 0;
    for (let i = 0; i < this.length; i++) {
        const char = this.charCodeAt(i);
        hash = (hash << 5) - hash + char;
        hash |= 0;
    }
    return hash;
};
let postRequest = (url, data, callback, final) => {
    fetch(url, {
        method: 'POST',
        body: JSON.stringify(data),
        headers: {
            "Content-type": "application/json; charset=UTF-8"
        }
    }).then((response) => {
        const statusCode = response.status;
        return response.text().then((text) => ({ text, statusCode }));
    }).then(({ text, statusCode }) => callback(text, statusCode)).finally(() => final());
};