const languages = {
    go: {
        name: 'GoLang',
        highlighter: 'go',
        renderer: 'codemirror',
        actions: 'run',
    },
    java: {
        name: 'Java',
        highlighter: 'text/x-java',
        renderer: 'codemirror',
        actions: 'run',
    },
    json: {
        name: 'JSON',
        highlighter: 'application/json',
        renderer: 'codemirror',
        actions: 'none',
    },
    markdown: {
        name: 'Markdown Edit',
        highlighter: 'text/x-markdown',
        renderer: 'codemirror',
        actions: 'view',
    },
    markdown_view: {
        name: 'Markdown View',
        highlighter: null,
        renderer: 'markdown',
        actions: 'edit',
    },
    mysql8: {
        name: 'MySQL 8',
        highlighter: 'sql',
        renderer: 'codemirror',
        actions: 'run',
    },
    php82: {
        name: 'PHP 8.2',
        highlighter: 'php',
        renderer: 'codemirror',
        actions: 'run',
    },
    postgres13: {
        name: 'PostgreSQL 13',
        highlighter: 'sql',
        renderer: 'codemirror',
        actions: 'run',
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

    if (app.actions !== languages[app.lang].actions) {
        if (languages[app.lang].actions === 'run') {
            editButton.style.display = 'none';
            viewButton.style.display = 'none';
            runButton.style.display = '';
            cleanResultButton.style.display = '';
        } else if (languages[app.lang].actions === 'view') {
            editButton.style.display = 'none';
            viewButton.style.display = '';
            runButton.style.display = 'none';
            cleanResultButton.style.display = 'none';
        } else if (languages[app.lang].actions === 'edit') {
            editButton.style.display = '';
            viewButton.style.display = 'none';
            runButton.style.display = 'none';
            cleanResultButton.style.display = 'none';
        } else { // none
            editButton.style.display = 'none';
            viewButton.style.display = 'none';
            runButton.style.display = 'none';
            cleanResultButton.style.display = 'none';
        }
        app.actions = languages[app.lang].actions;
    }

    langSelect.value = app.lang;
};
setLang(localStorage['initialLang']);
