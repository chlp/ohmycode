import {app} from "./app.js";
import {actions, onFileChange} from "./connect.js";
import {contentCodeMirror, contentCodeMirrorBlock, contentMarkdownBlock} from "./editor.js";

const languages = {
    go: {
        name: 'GoLang',
        highlighter: 'go',
        renderer: 'codemirror',
        actions: 'run',
        helloWorld: 'go',
    },
    java: {
        name: 'Java',
        highlighter: 'text/x-java',
        renderer: 'codemirror',
        actions: 'run',
        helloWorld: 'java',
    },
    json: {
        name: 'JSON',
        highlighter: 'application/json',
        renderer: 'codemirror',
        actions: 'none',
        helloWorld: 'json',
    },
    markdown: {
        name: 'Markdown Edit',
        highlighter: 'text/x-markdown',
        renderer: 'codemirror',
        actions: 'view',
        helloWorld: 'markdown',
    },
    markdown_view: {
        name: 'Markdown View',
        highlighter: null,
        renderer: 'markdown',
        actions: 'edit',
        helloWorld: undefined,
    },
    mysql8: {
        name: 'MySQL 8',
        highlighter: 'sql',
        renderer: 'codemirror',
        actions: 'run',
        helloWorld: 'mysql',
    },
    php82: {
        name: 'PHP 8.2',
        highlighter: 'php',
        renderer: 'codemirror',
        actions: 'run',
        helloWorld: 'php',
    },
    postgres13: {
        name: 'PostgreSQL 13',
        highlighter: 'sql',
        renderer: 'codemirror',
        actions: 'run',
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
langSelect.onchange = () => {
    setLang(langSelect.value);
};

const langChangeHandlers = [];
const onLangChange = (callback) => {
    if (typeof callback === "function") {
        langChangeHandlers.push(callback);
    }
};

let currentAction = undefined;

const setLang = (lang) => {
    if (app.lang === lang) {
        return;
    }

    if (languages[lang] === undefined) {
        lang = 'markdown';
    }
    app.lang = lang;

    contentCodeMirror.setOption('mode', languages[app.lang].highlighter);

    if (app.renderer !== languages[app.lang].renderer) {
        if (languages[app.lang].renderer === 'markdown') {
            contentCodeMirrorBlock.style.display = 'none';
            contentMarkdownBlock.style.display = '';
        } else { // codemirror for else
            contentCodeMirrorBlock.style.display = '';
            contentMarkdownBlock.style.display = 'none';
            contentCodeMirror.refresh()
        }
        app.renderer = languages[app.lang].renderer;
    }

    currentAction = languages[app.lang];

    langChangeHandlers.forEach(fn => fn(languages[app.lang]));

    langSelect.value = app.lang;

    if (typeof actions !== 'undefined') {
        actions.setLang(app.lang);
    }
    localStorage['initialLang'] = app.lang;
    contentCodeMirror.focus();
};

onFileChange(setLang);

const getLang = (langName) => {
    return languages[langName];
};

const getAction = () => {
    return currentAction;
};

export {getLang, setLang, onLangChange, getAction};