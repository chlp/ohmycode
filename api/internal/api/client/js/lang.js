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
