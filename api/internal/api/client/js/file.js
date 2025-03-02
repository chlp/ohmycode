const currentWriterInfo = document.getElementById('current-writer-info');
const runButton = document.getElementById('run-button');
const cleanResultButton = document.getElementById('clean-result-button');
const runnerContainerBlock = document.getElementById('runner-container');
const runnerEditButton = document.getElementById('runner-edit-button');
const runnerInput = document.getElementById('runner-input');
const runnerSaveButton = document.getElementById('runner-save-button');
const contentContainerBlock = document.getElementById('content-container');
const resultContainerBlock = document.getElementById('result-container');
const controlsContainerBlock = document.getElementById('controls-container');
const langSelect = document.getElementById('lang-select');

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
        // todo: show hint
        console.log('someone else is changing content now. wait please:', file.writer_id, app.id);
        return;
    }
    if (file.writer_id === '') {
        file.writer_id = app.id;
    }
});

contentCodeMirror.on('drop', (cm, event) => {
    event.preventDefault();
});

document.addEventListener('dragover', (event) => {
    event.preventDefault();
});

document.addEventListener('drop', (event) => {
    event.preventDefault();
    const droppedFiles = event.dataTransfer.files;
    if (droppedFiles.length === 0) {
        return;
    }
    const droppedFile = droppedFiles[0];
    if (droppedFile.size > 512 * 1024) {
        console.warn('File too large (>512Kb)', droppedFile);
        return;
    }
    const reader = new FileReader();
    reader.onload = async (e) => {
        if (await isFileBinary(droppedFile)) {
            console.warn("Wrong file (binary)", droppedFile);
            return;
        }
        if (file.writer_id !== '' && file.writer_id !== app.id) {
            return;
        }

        let newFileName = droppedFile.name;
        let newContent = e.target.result;
        const allowedCharsRegex = /[^0-9a-zA-Z_!?:=+\-,.\s'\u0400-\u04ff]/g;
        newFileName = newFileName.replace(allowedCharsRegex, '');
        newFileName = newFileName.substring(0, 64);

        fileNameBlock.innerHTML = newFileName;
        file.name = newFileName;
        actions.setFileName(newFileName);

        contentCodeMirror.setValue(newContent);
        contentMarkdownBlock.innerHTML = marked.parse(file.content);
        actions.setContent(newContent);
    };
    reader.onerror = function () {
        console.error('Error occurred: ' + droppedFile);
    };
    reader.readAsText(droppedFile);
});

let currentLang, currentRenderer;

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
    actions.setLang(currentLang);
    contentCodeMirror.focus();
};

const setLang = (newLang) => {
    if (currentLang === newLang) {
        return;
    }
    if (languages[newLang] === undefined) {
        newLang = 'markdown';
    }
    currentLang = newLang;
    contentCodeMirror.setOption('mode', languages[currentLang].highlighter);
    if (currentRenderer !== languages[currentLang].renderer) {
        if (languages[currentLang].renderer === 'markdown') {
            contentCodeMirrorBlock.style.display = 'none';
            contentMarkdownBlock.style.display = '';
        } else { // codemirror for else
            contentCodeMirrorBlock.style.display = '';
            contentMarkdownBlock.style.display = 'none';
            contentCodeMirror.refresh()
        }
        currentRenderer = languages[currentLang].renderer;
    }
    langSelect.value = currentLang;
};
setLang(localStorage['initialLang']);

let writerBlocksUpdate = () => {
    if (!app.isOnline) {
        contentCodeMirror.setOption('readOnly', true);
        currentWriterInfo.style.removeProperty('display');
        currentWriterInfo.innerHTML = 'Offline';
        return;
    }

    contentCodeMirror.setOption('readOnly', file.writer_id !== '' && file.writer_id !== app.id);
    if (file.writer_id === '' || file.writer_id === app.id) {
        currentWriterInfo.style.display = 'none';
        currentWriterInfo.innerHTML = '';
    } else {
        currentWriterInfo.style.removeProperty('display');
        currentWriterInfo.innerHTML = 'Editing is blocked by someone else';
    }
};

let runnerBlocksUpdate = () => {
    if (file.is_runner_online) {
        runnerContainerBlock.style.display = 'none';
    }
    runnerEditButton.style.display = file.is_runner_online ? 'none' : 'block';
};

let runnerEditButtonOnclick = () => {
    if (runnerContainerBlock.style.display === 'block') {
        runnerContainerBlock.style.display = 'none';
    } else {
        runnerContainerBlock.style.display = 'block';
        runnerInput.focus();
    }
};
runnerSaveButton.onclick = () => {
    actions.setRunner();
};
runnerInput.onkeydown = (event) => {
    if (event.key === 'Enter') {
        event.preventDefault();
        actions.setRunner();
    } else if (event.key === 'Escape') {
        runnerEditButtonOnclick();
    }
};

let socket = null;
let createWebSocket = () => {
    socket = new WebSocket(`${apiUrl}/file`);
    socket.onopen = () => {
        socket.send(JSON.stringify({
            action: 'init',
            file_id: file.id,
            app_id: app.id,
            user_id: app.userId,
            user_name: app.userName,
            lang: currentLang
        }));
    };
    socket.onclose = (event) => {
        if (event.wasClean) {
            console.log(`Connection closed, code=${event.code}, reason=${event.reason}`);
        } else {
            console.log('Connection closed with error');
        }
        app.isOnline = false;
        writerBlocksUpdate();
        socket = null;
    };
    socket.onerror = (error) => {
        console.log('WebSocket error: ', error);
    };
    socket.onmessage = (event) => {
        try {
            const data = JSON.parse(event.data);
            if (data.error !== undefined) {
                console.log('file::pageUpdater: getUpdate error', data);
                return;
            }

            let previousWriterId = file.writer_id;

            if (typeof data.content === 'undefined') {
                data.content = file.content;
            }
            file = data;

            if (file.persisted) {
                FilesHistory.saveFileToDB(file.id, file.name, file.content_updated_at);
            }
            document.title = `OhMyCode – ${file.name}`;

            // update users
            updateUsers();

            // update "Code is writing now by" block
            writerBlocksUpdate();

            // update runner ui
            runnerBlocksUpdate();

            // update result ui
            resultBlockUpdate();

            // update file name
            if (fileNameBlock.innerHTML !== file.name && !fileNameEditing) {
                fileNameBlock.innerHTML = file.name;
            }

            // update code
            if (
                !app.isOnline || // first load
                (
                    file.writer_id !== app.id && previousWriterId !== app.id && // do not update if current user is writer
                    ohMySimpleHash(file.content) !== ohMySimpleHash(contentCodeMirror.getValue()) // do not update if code is the same already
                )
            ) {
                let {left, top} = contentCodeMirror.getScrollInfo();
                let {line, ch} = contentCodeMirror.getCursor();
                contentCodeMirror.setValue(file.content);
                contentMarkdownBlock.innerHTML = marked.parse(file.content);
                contentCodeMirror.scrollTo(left, top);
                contentCodeMirror.setCursor({line: line, ch: ch});
            }

            // update lang
            setLang(file.lang);

            if (controlsContainerBlock.style.display !== 'block') {
                controlsContainerBlock.style.display = 'block';
            }

            if (!app.isOnline) {
                app.isOnline = true;
                contentSender();
                writerBlocksUpdate();
            }
        } catch (error) {
            console.error('Wrong message:', error);
        }
    };
};
createWebSocket();
let reconnectAttempts = 0;
setInterval(() => {
    if (socket === null) {
        reconnectAttempts++;
        createWebSocket();
    } else {
        reconnectAttempts = 0;
    }
}, 1000 * Math.min(2 ** reconnectAttempts, 30) + Math.random() * 3000);

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
