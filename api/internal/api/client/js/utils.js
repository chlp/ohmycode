const genUuid = () => { // Генерация случайного UUID без дефисов
    return 'xxxxxxxxxxxx4xxxyxxxxxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
        const r = Math.random() * 16 | 0, v = c === 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
};

const isUuid = (id) => {
    return (new RegExp(`^[a-z0-9]{32}$`)).test(id);
};

const getFileIdFromUrl = () => {
    const url = new URL(window.location.href);
    let fileId = url.pathname.replace(/^\/|\/$/g, '');
    if (!isUuid(fileId)) {
        fileId = genUuid();
        history.pushState({}, null, '/' + fileId);
    }
    return fileId;
}

const languages = {
    go: {
        name: 'GoLang',
        highlighter: 'go',
        renderer: 'codemirror',
    },
    java: {
        name: 'Java',
        highlighter: 'text/x-java',
        renderer: 'codemirror',
    },
    json: {
        name: 'JSON',
        highlighter: 'application/json',
        renderer: 'codemirror',
    },
    markdown: {
        name: 'Markdown',
        highlighter: 'text/x-markdown',
        renderer: 'codemirror',
    },
    markdown_view: {
        name: 'Markdown View',
        highlighter: null,
        renderer: 'markdown',
    },
    mysql8: {
        name: 'MySQL 8',
        highlighter: 'sql',
        renderer: 'codemirror',
    },
    php82: {
        name: 'PHP 8.2',
        highlighter: 'php',
        renderer: 'codemirror',
    },
    postgres13: {
        name: 'PostgreSQL 13',
        highlighter: 'sql',
        renderer: 'codemirror',
    },
};

let ohMySimpleHash = (str) => {
    if (str === undefined) {
        return 0;
    }
    let hash = 0;
    for (let i = 0; i < str.length; i++) {
        const char = str.charCodeAt(i);
        hash = (hash << 5) - hash + char;
        hash |= 0;
    }
    return hash;
};

let postRequest = (action, data, callback) => {
    try {
        socket.send(JSON.stringify({
            ...data,
            action: action,
        }));
    } finally {
        if (typeof callback === 'function') {
            callback();
        }
    }
};

let copyToClipboard = (text) => {
    if (navigator.clipboard && window.isSecureContext) {
        return navigator.clipboard.writeText(text).then(() => {
            console.log("Text copied to clipboard");
        }).catch(err => {
            console.error("Failed to copy: ", err);
        });
    } else {
        let textArea = document.createElement("textarea");
        textArea.value = text;

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
            document.body.removeChild(textArea);
        }
    }
};

let saveContentToFile = (text, fileName) => {
    const blob = new Blob([text], {type: 'text/plain'});
    const a = document.createElement('a');
    a.href = URL.createObjectURL(blob);
    if (!/\.[0-9a-z]+$/i.test(fileName)) {
        fileName += '.txt';
    }
    a.download = fileName;
    a.style.display = 'none';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(a.href);
};

document.addEventListener('keydown', function (event) {
    if ((event.ctrlKey || event.metaKey) && event.key === 's') {
        event.preventDefault();
        saveContentToFile();
    }
});

let isFileBinary = async (file) => {
    const buffer = await file.arrayBuffer();
    const bytes = new Uint8Array(buffer);
    const maxBytesToCheck = Math.min(bytes.length, 32 * 1024);
    let nonPrintableCount = 0;
    for (let i = 0; i < maxBytesToCheck; i++) {
        const byte = bytes[i];
        if ((byte < 32 || byte > 126) && byte !== 9 && byte !== 10 && byte !== 13) {
            nonPrintableCount++;
        }
    }
    let nonPrintableRateForBinary = 0.4;
    if (bytes.length < 1024) {
        nonPrintableRateForBinary = 0.6;
    }
    return nonPrintableCount / maxBytesToCheck > nonPrintableRateForBinary;
}
