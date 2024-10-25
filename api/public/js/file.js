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

let currentWriterInfo = document.getElementById('current-writer-info');
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


let filePreviousState = {};
let isOnline = true;

let userId = localStorage['userId'];
if (userId === undefined) {
    userId = genUuid();
    localStorage['userId'] = userId;
}

let userName = localStorage['user_name'];
if (userName === undefined) {
    userName = '';
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
    if (!isOnline) {
        contentBlock.setOption('readOnly', true);
        currentWriterInfo.style.removeProperty('display');
        currentWriterInfo.innerHTML = 'Offline';
        return;
    }

    contentBlock.setOption('readOnly', file.writer_id !== '' && file.writer_id !== userId);
    if (file.writer_id === '' || file.writer_id === userId) {
        currentWriterInfo.style.display = 'none';
        currentWriterInfo.innerHTML = '';
    } else {
        let writerName = '';
        if (file.users[file.writer_id]) {
            writerName = file.users[file.writer_id].name;
        } else {
            writerName = '???';
        }
        currentWriterInfo.style.removeProperty('display');
        currentWriterInfo.innerHTML = 'Code is writing now by ' + writerName;
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

let isResultFilledWithInProgress = false;
let resultBlockUpdate = () => {
    let isRunBtnShouldBeDisabled = false;
    if (file.is_waiting_for_result) {
        isRunBtnShouldBeDisabled = true;
        if (isResultFilledWithInProgress) {
            resultBlock.setValue(resultBlock.getValue() + '.');
        } else {
            isResultFilledWithInProgress = true;
            resultBlock.setValue('In progress...');
        }
    } else if (file.result.length > 0) {
        if (
            isResultFilledWithInProgress ||
            ohMySimpleHash(filePreviousState.result) !== ohMySimpleHash(file.result)
        ) {
            isResultFilledWithInProgress = false;
            resultBlock.setValue(file.result);
        }
    } else if (file.is_runner_online) {
        isResultFilledWithInProgress = false;
        resultBlock.setValue('runner will write result here...');
    } else {
        isRunBtnShouldBeDisabled = true;
        isResultFilledWithInProgress = false;
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
    }, (response, statusCode) => {
        response = response.trim();

        if (statusCode === 200 || statusCode === 204) {
            lastUpdateTimestamp = +new Date;
            if (!isOnline) {
                isOnline = true;
                writerBlocksUpdate();
            }
        }

        if (statusCode === 204 || response.length === 0) {
            resultBlockUpdate(); // adding more dots to "In progress..."
            return;
        }

        if (statusCode !== 200) {
            console.error("file::pageUpdater: error", statusCode, response)
            return;
        }

        let data = {};
        try {
            data = JSON.parse(response);
        } catch (error) {
            console.error("file::pageUpdater: failed to parse JSON:", error, response);
            return;
        }

        if (data.error !== undefined) {
            console.log('file::pageUpdater: getUpdate error', data);
            return;
        }

        filePreviousState = {...file};
        file = data;

        // update users
        updateUsers();

        // update "Code is writing now by" block
        writerBlocksUpdate();

        // update runner ui
        runnerBlocksUpdate();

        // update result ui
        resultBlockUpdate();

        // update file name
        if (filePreviousState.name !== file.name && !fileNameEditing) {
            fileNameBlock.innerHTML = file.name;
        }

        // update code
        if (
            file.writer_id !== userId && // do not update if current user is writer
            ohMySimpleHash(filePreviousState.code) !== ohMySimpleHash(file.content) // do not update if code is the same already
        ) {
            let {left, top} = contentBlock.getScrollInfo();
            let {line, ch} = contentBlock.getCursor();
            contentBlock.setValue(file.content);
            contentBlock.scrollTo(left, top);
            contentBlock.setCursor({line: line, ch: ch});
        }

        // update lang
        if (filePreviousState.lang !== file.lang) {
            currentLang = file.lang;
            langSelect.value = currentLang;
            contentBlock.setOption('mode', languages[currentLang].highlighter);
        }

        controlsContainerBlock.style.display = 'block';
    }, () => {
        if (isOnline && +new Date - lastUpdateTimestamp > 5000) {
            isOnline = false;
            writerBlocksUpdate();
        }

        console.log(isOnline ? 1000 : 5000)
        clearTimeout(pageUpdaterTimer);
        pageUpdaterTimer = setTimeout(() => {
            pageUpdater();
        }, isOnline ? 1000 : 5000);

        pageUpdaterIsInProgress = false;
    });
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
