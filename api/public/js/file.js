let fileId = window.location.pathname.slice(1);
if (!isUuid(fileId)) {
    fileId = genUuid();
    history.pushState({}, null, '/' + fileId);
}
let file = {
    "id": fileId,
    "name": "",
    "content": "",
    "lang": 'markdown',
    "runner": "",
    "is_runner_online": false,
    "updated_at": null,
    "writer_id": "",
    "users": [],
    "is_waiting_for_result": false,
    "result": ""
};

let sessionStatusBlock = document.getElementById('session-status');
let currentWriterInfo = document.getElementById('current-writer-info');
let currentWriterName = document.getElementById('current-writer-name');
let runButton = document.getElementById('run-button');
let cleanResultButton = document.getElementById('clean-result-button');
let runnerContainerBlock = document.getElementById('runner-container');
let runnerEditButton = document.getElementById('runner-edit-button');
let runnerInput = document.getElementById('runner-input');
let runnerSaveButton = document.getElementById('runner-save-button');
let codeContainerBlock = document.getElementById('code-container');
let resultContainerBlock = document.getElementById('result-container');
let controlsContainerBlock = document.getElementById('controls-container');
let langSelect = document.getElementById('lang-select');


let sessionPreviousState = {};
let sessionIsOnline = true;

let userId = localStorage['userId'];
if (userId === undefined) {
    userId = genUuid();
    localStorage['userId'] = userId;
}

let userName = localStorage['initialUserName'];
if (userName === undefined) {
    userName = randomName();
    localStorage['initialUserName'] = userName;
}

let currentLang = 'markdown';
if (localStorage['initialLang'] === undefined) {
    localStorage['initialLang'] = currentLang;
} else {
    currentLang = localStorage['initialLang'];
}
for (const key in languages) {
    if (languages.hasOwnProperty(key)) {
        const option = document.createElement('option');
        option.value = key;
        option.textContent = languages[key].name;
        langSelect.appendChild(option);
    }
}
langSelect.value = currentLang;

let getCodeTheme = () => {
    // https://codemirror.net/5/demo/theme.html
    // todo: temporary turned off light/dark scheme changing
    return 'base16-dark';
    // if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
    //     return 'base16-dark';
    // }
    // return 'base16-light';
};
let getResultTheme = () => {
    return 'tomorrow-night-bright';
};
let contentBlock = CodeMirror.fromTextArea(document.getElementById('code'), {
    lineNumbers: true,
    mode: languages[currentLang].highlighter, // javascript, go, php, sql
    matchBrackets: true,
    indentWithTabs: false,
    tabSize: 4,
    theme: getCodeTheme(),
    autofocus: true,
});
contentBlock.on('keydown', function (codemirror, event) {
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
    if (file.writer_id !== '' && file.writer_id !== userId) {
        // todo: show hint
        console.log('someone else is changing code now. wait please:', file.writer_id, userId);
    }
});
let resultBlock = CodeMirror.fromTextArea(document.getElementById('result'), {
    lineNumbers: true,
    readOnly: true,
    theme: getResultTheme(),
});
// todo: here we could show hint that it is not editable result
window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', event => {
    contentBlock.setOption('theme', getCodeTheme());
    resultBlock.setOption('theme', getResultTheme());
});

langSelect.onchange = () => {
    currentLang = langSelect.value;
    contentBlock.setOption('mode', languages[currentLang].highlighter);
    actions.setLang(currentLang);
};

let writerBlocksUpdate = () => {
    contentBlock.setOption('readOnly', file.writer_id !== '' && file.writer_id !== userId);
    let newWriterName = '?';
    if (file.writer_id === '' || file.writer_id === userId) {
        newWriterName = '';
        currentWriterInfo.style.display = 'none';
    } else {
        if (file.users[file.writer_id]) {
            newWriterName = file.users[file.writer_id].name;
        } else {
            newWriterName = '???';
        }
        currentWriterInfo.style.removeProperty('display');
    }
    currentWriterName.innerHTML = newWriterName;
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

let resultBlockUpdate = () => {
    let isRunBtnShouldBeDisabled = false;
    if (file.is_waiting_for_result) {
        isRunBtnShouldBeDisabled = true;
        if (resultBlock.getValue().startsWith('In progress')) {
            resultBlock.setValue(resultBlock.getValue() + '.');
        } else {
            resultBlock.setValue('In progress...');
        }
    } else if (file.result.length > 0) {
        if (ohMySimpleHash(sessionPreviousState.result) !== ohMySimpleHash(file.result)) {
            resultBlock.setValue(file.result);
        }
    } else if (file.is_runner_online) {
        resultBlock.setValue('runner will write result here...');
    } else {
        isRunBtnShouldBeDisabled = true;
        resultBlock.setValue('...');
    }

    if (isRunBtnShouldBeDisabled) {
        runButton.setAttribute('disabled', 'true');
    } else {
        runButton.removeAttribute('disabled');
    }

    if (file.is_waiting_for_result || file.result.length > 0) {
        resultContainerBlock.style.display = 'block';
        codeContainerBlock.style.height = 'calc(68vh - 90px)';
        cleanResultButton.removeAttribute('disabled');
    } else {
        resultContainerBlock.style.display = 'none';
        codeContainerBlock.style.height = 'calc(98vh - 90px)';
        cleanResultButton.setAttribute('disabled', 'true');
    }
};
resultBlockUpdate();
document.addEventListener('DOMContentLoaded', () => {
    setTimeout(() => {
        codeContainerBlock.style.transition = 'height 0.5s ease';
    }, 100);
});

let lastUpdateTimestamp = +new Date;
let pageUpdaterTimer = 0;
let pageUpdaterIsInProgress = false;
let pageUpdater = () => {
    let start = +new Date;
    if (pageUpdaterIsInProgress) {
        return;
    }
    pageUpdaterIsInProgress = true;
    postRequest('/file/get', {
        file_id: fileId,
        user_id: userId,
        user_name: userName,
        lang: currentLang,
        last_update: file.updated_at,
        is_keep_alive: true,
    }, (response) => {
        response = response.trim();
        pageUpdaterIsInProgress = false;
        lastUpdateTimestamp = +new Date;
        if (response.length === 0) {
            resultBlockUpdate(); // adding more dots to "In progress..."
            return;
        }
        let data = {};
        try {
            data = JSON.parse(response);
        } catch (error) {
            console.error("session::pageUpdater: failed to parse JSON:", error, response);
            return;
        }

        if (data.error !== undefined) {
            console.log('session::pageUpdater: getUpdate error', data);
            return;
        }

        sessionPreviousState = {...file};
        file = data;

        // update users
        updateUsers();

        // update "Code is writing now by" block
        writerBlocksUpdate();

        // update runner ui
        runnerBlocksUpdate();

        // update result ui
        resultBlockUpdate();

        // update session name
        if (sessionPreviousState.name !== file.name && !sessionNameEditing) {
            sessionNameBlock.innerHTML = file.name;
        }

        // update code
        if (
            file.writer_id !== userId && // do not update if current user is writer
            ohMySimpleHash(sessionPreviousState.code) !== ohMySimpleHash(file.content) // do not update if code is the same already
        ) {
            let {left, top} = contentBlock.getScrollInfo();
            let {line, ch} = contentBlock.getCursor();
            contentBlock.setValue(file.content);
            contentBlock.scrollTo(left, top);
            contentBlock.setCursor({line: line, ch: ch});
        }

        // update lang
        if (sessionPreviousState.lang !== file.lang) {
            currentLang = file.lang;
            langSelect.value = currentLang;
            contentBlock.setOption('mode', languages[currentLang].highlighter);
        }

        controlsContainerBlock.style.display = 'block';
    }, () => {
        clearTimeout(pageUpdaterTimer);
        pageUpdaterTimer = setTimeout(() => {
            pageUpdater();
        }, 1000);
    });
    if (+new Date - lastUpdateTimestamp > 10000) { // more than 10 seconds
        if (sessionIsOnline) {
            sessionIsOnline = false;
            sessionStatusBlock.classList.remove('online');
            sessionStatusBlock.classList.add('offline');
            sessionStatusBlock.innerHTML = ' offline';
        }
    } else {
        if (!sessionIsOnline) {
            sessionIsOnline = true;
            sessionStatusBlock.classList.remove('offline');
            sessionStatusBlock.classList.add('online');
            sessionStatusBlock.innerHTML = '';
        }
    }
};
pageUpdater();

let codeSenderTimer = 0;
let codeSender = () => {
    if (ohMySimpleHash(file.content) !== ohMySimpleHash(contentBlock.getValue())) {
        actions.setCode(() => {
            clearTimeout(codeSenderTimer);
            codeSenderTimer = setTimeout(() => {
                codeSender();
            }, 1000);
        });
    } else {
        clearTimeout(codeSenderTimer);
        codeSenderTimer = setTimeout(() => {
            codeSender();
        }, 300);
    }
};
codeSender();

let runCode = () => {
    if (!file.is_runner_online) {
        resultBlock.setValue('No runner is available to run your code :(');
        return;
    }
    clearTimeout(pageUpdaterTimer);
    let runCodeCall = () => {
        file.result = 'In progress..';
        resultBlock.setValue('In progress..');
        runButton.setAttribute('disabled', 'true');
        actions.runCode(pageUpdater);
    };
    if (ohMySimpleHash(file.content) !== ohMySimpleHash(contentBlock.getValue())) {
        actions.setCode(() => {
            runCodeCall();
        });
    } else {
        runCodeCall();
    }
};
runButton.onclick = () => {
    runCode();
};
codeContainerBlock.onkeydown = (event) => {
    if ((event.metaKey || event.ctrlKey) && event.key === 'Enter') {
        runCode();
    }
};

cleanResultButton.onclick = () => {
    file.result = '';
    resultBlock.setValue('');
    actions.cleanResult(() => {
        resultBlockUpdate();
    });
};
