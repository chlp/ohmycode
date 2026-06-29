import {ohMySimpleHash} from "./utils.js";
import {app, file} from "./app.js";
import {actions} from "./connect.js";
import {getCurrentLang, onLangChange, setLang} from "./lang.js";
import {historyPanelToggle} from "./sidebar.js";
import {setLockStatus, setIdleStatus} from "./status.js";

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

const contentMaxLengthKb = 512;

const lockButton = document.getElementById('header-lock-btn');
const lockIconClosed = `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path></svg>`;
const lockIconOpen = `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect><path d="M7 11V7a5 5 0 0 1 9.9-1"></path></svg>`;

let updateEditorLockStatus = () => {
    if (app.isROLink) {
        contentCodeMirror.setOption('readOnly', true);
        setLockStatus('Read-only');
        lockButton.style.display = '';
        lockButton.innerHTML = lockIconClosed;
        lockButton.title = 'Read-only access';
        lockButton.disabled = true;
        return;
    }

    lockButton.disabled = false;
    lockButton.style.display = '';

    if (!app.isOnline) {
        contentCodeMirror.setOption('readOnly', true);
        setLockStatus('Offline');
        lockButton.innerHTML = lockIconOpen;
        lockButton.title = 'Lock editing';
        return;
    }

    if (file.is_locked) {
        contentCodeMirror.setOption('readOnly', true);
        setLockStatus('Locked');
        lockButton.innerHTML = lockIconClosed;
        lockButton.title = 'Unlock editing';
        return;
    }

    lockButton.innerHTML = lockIconOpen;
    lockButton.title = 'Lock editing';
    contentCodeMirror.setOption('readOnly', file.writer_id !== '' && file.writer_id !== app.id);
    if (file.writer_id === '' || file.writer_id === app.id) {
        setLockStatus(null);
    } else {
        setLockStatus('Editing is blocked by someone else');
    }
};

const updateContentSizeStatus = () => {
    const sizeKb = (contentCodeMirror.getValue().length / 1024).toFixed(1);
    setIdleStatus(`${sizeKb} KB / ${contentMaxLengthKb} KB`);
};
contentCodeMirror.on('change', updateContentSizeStatus);

let contentSenderTimer = 0;
let contentSender = () => {
    let getNextUpdateFunc = (timeout) => {
        clearTimeout(contentSenderTimer);
        contentSenderTimer = setTimeout(() => {
            contentSender();
        }, timeout);
    };
    const newContent = contentCodeMirror.getValue();
    if (app.isOnline && !file.is_locked && ohMySimpleHash(file.content) !== ohMySimpleHash(newContent)) {
        getNextUpdateFunc(1000);
        contentMarkdownBlock.innerHTML = marked.parse(newContent); // todo: should have function to update all editors/views or load data after changing mode
        actions.setContent(newContent);
    } else {
        getNextUpdateFunc(500);
    }
};
setTimeout(() => {
    contentSender();
}, 500);

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

document.getElementById('header-download-btn').onclick = () => {
    saveContentToFile(file.content, file.name);
};

document.addEventListener('keydown', function (event) {
    if ((event.ctrlKey || event.metaKey) && event.key === 's') {
        event.preventDefault();
        saveContentToFile(file.content, file.name);
    }
});

document.onkeydown = (event) => {
    if ((event.metaKey || event.ctrlKey) && event.key === 'Enter') {
        const action = getCurrentLang().action;
        switch (action) {
            case 'run':
                actions.runTaskWithContent(contentCodeMirror.getValue());
                break;
            case 'view':
                setLang('markdown_view'); // todo: not only markdown
                break;
            case 'edit':
                setLang('markdown');
                break;
            default:
                console.warn("wrong action: ", action);
        }
    } else if (event.key === 'Escape') {
        historyPanelToggle();
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

onLangChange((lang) => {
    switch (lang.action) {
        case 'run':
            editButton.style.display = 'none';
            viewButton.style.display = 'none';
            break;
        case 'view':
            editButton.style.display = 'none';
            viewButton.style.display = '';
            break;
        case 'edit':
            editButton.style.display = '';
            viewButton.style.display = 'none';
            break;
        case 'none':
        default:
            editButton.style.display = 'none';
            viewButton.style.display = 'none';
            break;
    }
});

lockButton.onclick = () => {
    actions.setLocked(!file.is_locked);
};

export {contentCodeMirror, contentCodeMirrorBlock, contentMarkdownBlock, updateEditorLockStatus};