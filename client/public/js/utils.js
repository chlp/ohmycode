let genUuid = () => { // Генерация случайного UUID без дефисов
    return 'xxxxxxxxxxxx4xxxyxxxxxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
        const r = Math.random() * 16 | 0, v = c === 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
};

let isUuid = (id) => {
    return (new RegExp(`^[a-z0-9]{32}$`)).test(id);
};

const languages = {
    go: {
        name: 'GoLang',
        highlighter: 'go',
    },
    java: {
        name: 'Java',
        highlighter: 'text/x-java',
    },
    json: {
        name: 'JSON',
        highlighter: 'application/json',
    },
    markdown: {
        name: 'Markdown',
        highlighter: 'text/x-markdown',
    },
    mysql8: {
        name: 'MySQL 8',
        highlighter: 'sql',
    },
    php82: {
        name: 'PHP 8.2',
        highlighter: 'php',
    },
    postgres13: {
        name: 'PostgreSQL 13',
        highlighter: 'sql',
    }
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

let saveContentToFile = () => {
    const text = contentBlock.getValue();
    const blob = new Blob([text], { type: 'text/plain' });
    const a = document.createElement('a');
    a.href = URL.createObjectURL(blob);
    a.download = file.name + '.txt';
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


let processId = genUuid();
let checkForMultipleTabs = () => {
    let statusIdKey = 'file-status-id-' + fileId;
    let statusUpdatedAtKey = 'file-status-updatedAt-' + fileId;
    if (
        localStorage[statusIdKey] === undefined ||
        localStorage[statusIdKey] === processId ||
        +new Date - localStorage[statusUpdatedAtKey] > 2000
    ) {
        localStorage[statusIdKey] = processId;
        localStorage[statusUpdatedAtKey] = +new Date;
        return;
    }

    // stopping all intervals and timers and ask to close window
    let newTimerId = setTimeout(() => {
    }, 1);
    for (let i = 0; i <= newTimerId; i++) {
        clearTimeout(i);
    }
    let newIntervalId = setInterval(() => {
    }, 1);
    for (let i = 0; i <= newIntervalId; i++) {
        clearInterval(i);
    }

    document.title = '! OhMyCode';
    setInterval(() => {
        document.title = '! OhMyCode';
        setTimeout(() => {
            document.title = '? OhMyCode';
        }, 1000);
    }, 2000);
    document.body.innerHTML = '<h1 style="text-align: center; margin-top: 2em;">OhMyCode cannot work with one shared session in multiple tabs.<br>Please use only one tab for one session in one browser.</h1>';
};
setInterval(() => {
    checkForMultipleTabs();
}, 2000);