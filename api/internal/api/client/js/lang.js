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
    actions.setLang(app.lang);
    localStorage['initialLang'] = app.lang;
    contentCodeMirror.focus();
};

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
    langSelect.value = app.lang;
};
setLang(localStorage['initialLang']);
