import {app} from "./app.js";
import {actions, onFileChange} from "./connect.js";
import {contentCodeMirror, contentCodeMirrorBlock, contentMarkdownBlock} from "./editor.js";

const languages = {
    go: {
        key: 'go',
        name: 'GoLang',
        highlighter: 'go',
        renderer: 'codemirror',
        action: 'run',
        helloWorld: 'go',
    },
    java: {
        key: 'java',
        name: 'Java',
        highlighter: 'text/x-java',
        renderer: 'codemirror',
        action: 'run',
        helloWorld: 'java',
    },
    json: {
        key: 'json',
        name: 'JSON',
        highlighter: 'application/json',
        renderer: 'codemirror',
        action: 'none',
        helloWorld: 'json',
    },
    markdown: {
        key: 'markdown',
        name: 'Markdown Edit',
        highlighter: 'text/x-markdown',
        renderer: 'codemirror',
        action: 'view',
        helloWorld: 'markdown',
    },
    markdown_view: {
        key: 'markdown_view',
        name: 'Markdown View',
        highlighter: null,
        renderer: 'markdown',
        action: 'edit',
        helloWorld: undefined,
    },
    mysql8: {
        key: 'mysql8',
        name: 'MySQL 8',
        highlighter: 'sql',
        renderer: 'codemirror',
        action: 'run',
        helloWorld: 'mysql',
    },
    php82: {
        key: 'php82',
        name: 'PHP 8.2',
        highlighter: 'php',
        renderer: 'codemirror',
        action: 'run',
        helloWorld: 'php',
    },
    postgres13: {
        key: 'postgres13',
        name: 'PostgreSQL 13',
        highlighter: 'sql',
        renderer: 'codemirror',
        action: 'run',
        helloWorld: 'postgres',
    },
};

const langSelect = document.getElementById('lang-select');

for (const key in languages) {
    if (languages.hasOwnProperty(key)) {
        const option = document.createElement('option');
        option.value = key;
        option.textContent = languages[key].name;
        langSelect.appendChild(option);
    }
}

const langChangeHandlers = [];
const onLangChange = (callback) => {
    if (typeof callback === "function") {
        langChangeHandlers.push(callback);
    }
};

let currentLang = undefined;

const setLang = (langName) => {
    if (typeof currentLang !== 'undefined' && currentLang.lang === langName) {
        return;
    }

    if (typeof actions !== 'undefined') {
        actions.setLang(langName);
    }

    if (languages[langName] === undefined) {
        langName = 'markdown';
    }

    const langObj = languages[langName];
    const previousLang = currentLang;
    currentLang = langObj;

    contentCodeMirror.setOption('mode', langObj.highlighter);

    if (app.renderer !== langObj.renderer) {
        if (langObj.renderer === 'markdown') {
            contentCodeMirrorBlock.style.display = 'none';
            contentMarkdownBlock.style.display = '';
        } else { // codemirror for else
            contentCodeMirrorBlock.style.display = '';
            contentMarkdownBlock.style.display = 'none';
            contentCodeMirror.refresh()
        }
        app.renderer = langObj.renderer;
    }

    langSelect.onchange = () => {};
    langSelect.value = langName;
    langSelect.onchange = (ev) => {
        const changeToLangName = ev.target.value;
        localStorage['initialLang'] = changeToLangName;
        contentCodeMirror.focus();
        setLang(changeToLangName);
    };

    langChangeHandlers.forEach(fn => fn(langObj));
};

window.addEventListener("DOMContentLoaded", () => {
    onFileChange((file) => {
        setLang(file.lang);
    });
});

const getCurrentLang = () => {
    if (typeof currentLang === 'undefined') {
        return languages['markdown'];
    }
    return currentLang;
};

export {getCurrentLang, setLang, onLangChange};