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

let copyToClipboard = (text) => {
    if (navigator.clipboard && window.isSecureContext) {
        // Используем Clipboard API, если доступен
        return navigator.clipboard.writeText(text).then(() => {
            console.log("Text copied to clipboard");
        }).catch(err => {
            console.error("Failed to copy: ", err);
        });
    } else {
        // Fallback для старых браузеров
        let textArea = document.createElement("textarea");
        textArea.value = text;

        // Избегаем отображения элемента в окне
        textArea.style.position = "fixed";
        textArea.style.left = "-999999px";
        document.body.appendChild(textArea);
        textArea.focus();
        textArea.select();

        try {
            document.execCommand('copy');
            console.log("Text copied to clipboard");
        } catch (err) {
            console.error("Failed to copy: ", err);
        } finally {
            // Удаляем текстовое поле после копирования
            document.body.removeChild(textArea);
        }
    }
};

document.addEventListener('keydown', function(event) {
    if ((event.ctrlKey || event.metaKey) && event.key === 's') {
        event.preventDefault();
        console.log('Already saved :)');
    }
});