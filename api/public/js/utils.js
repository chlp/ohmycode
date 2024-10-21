let genUuid = () => {
    return crypto.randomUUID().replace(/-/g, '');
};

let isUuid = (id) => {
    return (new RegExp(`^[a-z0-9]{32}$`)).test(id);
};

let randomName = () => {
    const adjectives = [
        'Happy', 'Cheerful', 'Playful', 'Friendly', 'Bubbly', 'Jolly',
        'Witty', 'Quirky', 'Silly', 'Goofy', 'Sunny', 'Joyful',
        'Clever', 'Bouncy', 'Zippy', 'Peppy', 'Snazzy', 'Perky'
    ];
    const animals = [
        'Penguin', 'Panda', 'Koala', 'Bunny', 'Squirrel', 'Dolphin',
        'Turtle', 'Owl', 'Duckling', 'Kitten', 'Puppy', 'Sloth',
        'Raccoon', 'Goldfish', 'Hedgehog', 'Llama', 'Frog', 'Otter'
    ];

    const adjective = adjectives[Math.floor(Math.random() * adjectives.length)];
    const animal = animals[Math.floor(Math.random() * animals.length)];

    return `${adjective} ${animal}`;
}

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
    let hash = 0;
    for (let i = 0; i < str.length; i++) {
        const char = str.charCodeAt(i);
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
        return response.text().then((text) => ({text, statusCode}));
    }).then(({text, statusCode}) => callback(text, statusCode)).finally(() => final());
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

document.addEventListener('keydown', function (event) {
    if ((event.ctrlKey || event.metaKey) && event.key === 's') {
        event.preventDefault();
        console.log('Already saved :)');
    }
});