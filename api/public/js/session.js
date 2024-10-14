let sessionPreviousState = {...session};
sessionPreviousState.writer = '-'; // hack to init users
let sessionIsOnline = true;
let ping = undefined;
let userId = localStorage['userId'];
if (userId === undefined) {
    userId = initialUserId;
    localStorage['userId'] = userId;
}
let userName = session.users[userId] ? session.users[userId].name : undefined;
if (userName === undefined) {
    userName = localStorage['initialUserName'];
    if (userName === undefined) {
        userName = initialName;
        localStorage['initialUserName'] = userName;
    }
}

let checkForMultipleTabs = (isInitial) => {
    // todo: remove old sessions data
    let sessionStatusIdKey = 'session-status-id-' + session.id;
    let sessionStatusUpdatedAtKey = 'session-status-updatedAt-' + session.id;
    if (isInitial) {
        localStorage[sessionStatusIdKey] = initialUserId;
        localStorage[sessionStatusUpdatedAtKey] = +new Date;
    } else {
        if (localStorage[sessionStatusIdKey] !== initialUserId &&
            +new Date - localStorage[sessionStatusUpdatedAtKey] < 1500) {
            // stopping all intervals and timers and ask to close window
            let newTimerId = setTimeout(() => {
            });
            for (let i = 0; i < newTimerId; i++) {
                clearTimeout(i);
            }
            let newIntervalId = setInterval(() => {
            });
            for (let i = 0; i < newIntervalId; i++) {
                clearInterval(i);
            }
            document.title = '! OhMyCode';
            setInterval(() => {
                document.title = '! OhMyCode';
                setTimeout(() => {
                    document.title = '? OhMyCode';
                }, 1000);
            }, 2000);
            document.body.innerHTML = '<h1 style="text-align: center; margin-top: 2em;">OhMyCode cannot work with one session in multiple tabs.<br>Please leave just one tab for this session.</h1>';
        } else {
            localStorage[sessionStatusIdKey] = initialUserId;
            localStorage[sessionStatusUpdatedAtKey] = +new Date;
        }
    }
};
checkForMultipleTabs(true);
setInterval(() => {
    checkForMultipleTabs(false);
}, 2000);


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
let codeBlock = CodeMirror.fromTextArea(document.getElementById('code'), {
    lineNumbers: true,
    mode: langKeyToHighlighter[session.lang], // javascript, go, php, sql
    matchBrackets: true,
    indentWithTabs: false,
    tabSize: 4,
    theme: getCodeTheme(),
    autofocus: true,
});
codeBlock.on('keydown', function (codemirror, event) {
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
    if (session.writer !== '' && session.writer !== userId) {
        // todo: show hint
        console.log('someone else is changing code now. wait please:', session.writer, userId);
    }
});
let resultBlock = CodeMirror.fromTextArea(document.getElementById('result'), {
    lineNumbers: true,
    readOnly: true,
    theme: getResultTheme(),
});
// todo: here we could show hint that it is not editable result
window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', event => {
    codeBlock.setOption('theme', getCodeTheme());
    resultBlock.setOption('theme', getResultTheme());
});

let sessionStatusBlock = document.getElementById('session-status');
let currentWriterInfo = document.getElementById('current-writer-info');
let currentWriterName = document.getElementById('current-writer-name');
let langSelect = document.getElementById('lang-select');
let runButton = document.getElementById('run-button');
let cleanResultButton = document.getElementById('clean-result-button');
let runnerContainerBlock = document.getElementById('runner-container');
let runnerEditButton = document.getElementById('runner-edit-button');
let runnerInput = document.getElementById('runner-input');
let runnerSaveButton = document.getElementById('runner-save-button');
let codeContainerBlock = document.getElementById('code-container');
let resultContainerBlock = document.getElementById('result-container');

langSelect.onchange = () => {
    codeBlock.setOption('mode', langKeyToHighlighter[langSelect.value]);
    actions.setLang();
};

let writerBlocksUpdate = () => {
    codeBlock.setOption('readOnly', session.writer !== '' && session.writer !== userId);
    let newWriterName = '?';
    if (session.writer === '' || session.writer === userId) {
        newWriterName = '';
        currentWriterInfo.style.display = 'none';
    } else {
        if (session.users[session.writer]) {
            newWriterName = session.users[session.writer].name;
        } else {
            newWriterName = '???';
        }
        currentWriterInfo.style.removeProperty('display');
    }
    currentWriterName.innerHTML = newWriterName;
};
document.addEventListener('DOMContentLoaded', () => {
    writerBlocksUpdate();
});

let runnerBlocksUpdate = () => {
    if (session.runnerIsOnline) {
        runnerContainerBlock.style.display = 'none';
    }
    runnerEditButton.style.display = session.runnerIsOnline ? 'none' : 'block';
};
document.addEventListener('DOMContentLoaded', () => {
    runnerBlocksUpdate();
});

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
    if (session.isWaitingForResult) {
        if (resultBlock.getValue().startsWith('In progress')) {
            resultBlock.setValue(resultBlock.getValue() + '.');
        } else {
            resultBlock.setValue('In progress...');
        }
    } else if (session.result.length > 0) {
        if (sessionPreviousState.result.ohMySimpleHash() !== session.result.ohMySimpleHash()) {
            resultBlock.setValue(session.result);
        }
    } else if (session.runnerIsOnline) {
        resultBlock.setValue('runner will write result here...');
    } else {
        resultBlock.setValue('...');
    }
    if (session.isWaitingForResult || session.result.length > 0) {
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
    resultBlockUpdate();
    setTimeout(() => {
        codeContainerBlock.style.transition = 'height 0.5s ease';
    }, 100);
});

let isDebug = false;
let lastUpdateTimestamp = +new Date;
let pageUpdaterTimer = 0;
let pageUpdaterIsInProgress = false;
let pageUpdater = () => {
    let start = +new Date;
    if (pageUpdaterIsInProgress) {
        return;
    }
    pageUpdaterIsInProgress = true;
    postRequest('/action/session.php?action=getUpdate', {
        session: session.id,
        user: userId,
        userName: userName,
        lastUpdate: session.updatedAt ? session.updatedAt.date : null,
        action: 'getUpdate',
        isKeepAlive: true,
    }, (response) => {
        response = response.trim();
        pageUpdaterIsInProgress = false;
        ping = +new Date - start;
        if (isDebug) {
            console.log((new Date).toLocaleString() + ' | ping: ' + ping);
        }
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

        isNewSession = false;
        sessionPreviousState = {...session};
        session = data;

        // update users
        updateUsers();

        // update "Code is writing now by" block
        writerBlocksUpdate();

        // update runner ui
        runnerBlocksUpdate();

        // update result ui
        resultBlockUpdate();

        // update session name
        if (sessionPreviousState.name !== session.name && !sessionNameEditing) {
            sessionNameBlock.innerHTML = session.name;
        }

        // update code
        if (
            session.writer !== userId && // do not update if current user is writer
            sessionPreviousState.code.ohMySimpleHash() !== session.code.ohMySimpleHash() // do not update if code is the same already
        ) {
            let {left, top} = codeBlock.getScrollInfo();
            let {line, ch} = codeBlock.getCursor();
            codeBlock.setValue(session.code);
            codeBlock.scrollTo(left, top);
            codeBlock.setCursor({line: line, ch: ch});
        }

        // update lang
        if (sessionPreviousState.lang !== session.lang) {
            langSelect.value = session.lang;
            codeBlock.setOption('mode', langKeyToHighlighter[session.lang]);
        }
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
    if (session.code.ohMySimpleHash() !== codeBlock.getValue().ohMySimpleHash()) {
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
    if (!session.runnerIsOnline) {
        resultBlock.setValue('No runner is available to run your code :(');
        return;
    }
    clearTimeout(pageUpdaterTimer);
    if (session.code.ohMySimpleHash() !== codeBlock.getValue().ohMySimpleHash()) {
        actions.setCode(() => {
            actions.runCode(pageUpdater);
        });
    } else {
        actions.runCode(pageUpdater);
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
    actions.cleanCode(() => {
        resultBlockUpdate();
    });
};
