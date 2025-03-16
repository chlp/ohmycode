const contentContainerBlock = document.getElementById('content-container');
const contentMarkdownBlock = document.getElementById('content-markdown');

let contentCodeMirror = CodeMirror.fromTextArea(document.getElementById('content'), {
    lineNumbers: true,
    lineWrapping: true,
    readOnly: true,
    matchBrackets: true,
    indentWithTabs: false,
    tabSize: 4,
    theme: 'base16-dark',
    autofocus: true,
});
const contentCodeMirrorBlock = contentContainerBlock.getElementsByClassName('CodeMirror')[0];
contentCodeMirror.on('keydown', function (codemirror, event) {
    if ((event.ctrlKey || event.metaKey) && event.key === 'c') {
        return;
    }
    const nonTextKeys = [
        'Shift', 'Control', 'Alt', 'Meta', 'CapsLock', 'Tab',
        'Escape', 'ArrowUp', 'ArrowDown', 'ArrowLeft', 'ArrowRight',
        'Enter', 'Backspace', 'Delete', 'Home', 'End', 'PageUp',
        'PageDown', 'F1', 'F2', 'F3', 'F4', 'F5', 'F6', 'F7',
        'F8', 'F9', 'F10', 'F11', 'F12', 'Insert', 'Pause',
        'NumLock', 'ScrollLock', 'ContextMenu'
    ];
    if (nonTextKeys.includes(event.key)) {
        return;
    }
    if (file.writer_id !== '' && file.writer_id !== app.id) {
        console.log('someone else is changing content now. wait please:', file.writer_id, app.id);
        return;
    }
    if (file.writer_id === '') {
        file.writer_id = app.id;
    }
});

const statusBarBlock = document.getElementById('status-bar');
let updateEditorLockStatus = () => {
    if (!app.isOnline) {
        contentCodeMirror.setOption('readOnly', true);
        statusBarBlock.style.removeProperty('display');
        statusBarBlock.innerHTML = 'Offline';
        return;
    }

    contentCodeMirror.setOption('readOnly', file.writer_id !== '' && file.writer_id !== app.id);
    if (file.writer_id === '' || file.writer_id === app.id) {
        statusBarBlock.style.display = 'none';
        statusBarBlock.innerHTML = '';
    } else {
        statusBarBlock.style.removeProperty('display');
        statusBarBlock.innerHTML = 'Editing is blocked by someone else';
    }
};

let contentSenderTimer = 0;
let contentSender = () => {
    if (!app.isOnline) {
        return;
    }
    let getNextUpdateFunc = (timeout) => {
        clearTimeout(contentSenderTimer);
        contentSenderTimer = setTimeout(() => {
            contentSender();
        }, timeout);
    };
    const newContent = contentCodeMirror.getValue();
    if (ohMySimpleHash(file.content) !== ohMySimpleHash(newContent)) {
        contentMarkdownBlock.innerHTML = marked.parse(newContent);
        actions.setContent(newContent);
        getNextUpdateFunc(1000);
    } else {
        getNextUpdateFunc(500);
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

document.onkeydown = (event) => {
    if ((event.metaKey || event.ctrlKey) && event.key === 'Enter') {
        if (app.actions === 'run') {
            actions.runTask();
        } else if (app.actions === 'view') {
            setLang('markdown_view'); // todo: not only markdown
        } else if (app.actions === 'edit') {
            setLang('markdown');
        }
    }
};

const editButton = document.getElementById('edit-button');
const viewButton = document.getElementById('view-button');

viewButton.onclick = () => {
    setLang('markdown_view'); // todo: not only markdown
};

editButton.onclick = () => {
    setLang('markdown');
};
