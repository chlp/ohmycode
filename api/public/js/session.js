let sessionPreviousState = {...session};
sessionPreviousState.writer = '-'; // hack to init users
let sessionIsOnline = true;
let ping = undefined;
let isWriter = false;
let userId = localStorage['userId'];
if (userId === undefined) {
    userId = initialUserId;
    localStorage['userId'] = userId;
}
let userName = undefined;
session.users.forEach((user) => {
    if (user.id === userId) {
        userName = user.name;
    }
});
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
    if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
        return 'base16-dark';
    }
    return 'base16-light';
}
let codeBlock = CodeMirror.fromTextArea(document.getElementById('code'), {
    lineNumbers: true,
    mode: langKeyToHighlighter[session.lang], // javascript, go, php, sql
    matchBrackets: true,
    indentWithTabs: false,
    tabSize: 4,
    theme: getCodeTheme()
});
let resultBlock = CodeMirror.fromTextArea(document.getElementById('result'), {
    lineNumbers: true,
    readOnly: true,
    theme: getCodeTheme()
});
window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', event => {
    codeBlock.setOption('theme', getCodeTheme());
    resultBlock.setOption('theme', getCodeTheme());
});

let sessionNameContainerBlock = document.getElementById('session-name-container');
let sessionNameBlock = document.getElementById('session-name');
let sessionNameInput = document.getElementById('session-name-input');
let sessionNameSaveButton = document.getElementById('session-name-save-button');
let usersContainerBlock = document.getElementById('users-container');
let userNameContainerBlock = document.getElementById('user-name-container');
let userNameInput = document.getElementById('user-name-input');
let userNameSaveButton = document.getElementById('user-name-save-button');
let sessionStatusBlock = document.getElementById('session-status');
let becomeWriterButton = document.getElementById('become-writer-button');
let langSelect = document.getElementById('lang-select');
let runButton = document.getElementById('run-button');
let runnerContainerBlock = document.getElementById('runner-container');
let runnerEditButton = document.getElementById('runner-edit-button');
let runnerInput = document.getElementById('runner-input');
let runnerSaveButton = document.getElementById('runner-save-button');
let codeContainerBlock = document.getElementById('code-container');
let resultContainerBlock = document.getElementById('result-container');
resultContainerBlock.style.display = 'none';

sessionNameSaveButton.onclick = () => {
    actions.setSessionName()
};
sessionNameInput.onkeydown = (event) => {
    if (event.key === 'Enter') {
        event.preventDefault();
        actions.setSessionName();
    } else if (event.key === 'Escape') {
        sessionNameBlock.click();
    }
};
sessionNameBlock.onclick = () => {
    if (sessionNameContainerBlock.style.display === 'block') {
        sessionNameContainerBlock.style.display = 'none';
    } else {
        sessionNameInput.value = session.name;
        sessionNameContainerBlock.style.display = 'block';
        sessionNameInput.focus();
    }
};

userNameSaveButton.onclick = () => {
    actions.setUserName();
};
userNameInput.onkeydown = (event) => {
    if (event.key === 'Enter') {
        event.preventDefault();
        actions.setUserName();
    } else if (event.key === 'Escape') {
        ownUserNameOnclick();
    }
};
let ownUserNameOnclick = () => {
    if (userNameContainerBlock.style.display === 'block') {
        userNameContainerBlock.style.display = 'none';
    } else {
        userNameInput.value = userName;
        userNameContainerBlock.style.display = 'block';
        userNameInput.focus();
    }
};
langSelect.onchange = () => {
    if (isWriter) {
        codeBlock.setOption('mode', langKeyToHighlighter[langSelect.value]);
        actions.setLang();
    }
};

runButton.onclick = () => {
    actions.setRequest();
};
codeContainerBlock.onkeydown = (event) => {
    if ((event.metaKey || event.ctrlKey) && event.key === 'Enter') {
        actions.setRequest();
    }
};

let updateUsers = () => {
    if (sessionPreviousState.writer + JSON.stringify(sessionPreviousState.users) === session.writer + JSON.stringify(session.users)) {
        return;
    }
    let spectators = [];
    let writer = undefined;
    if (isNewSession) {
        isWriter = true;
        writer = {
            id: userId,
            name: userName,
            own: true,
        }
        spectators = []
    } else {
        isWriter = userId === session.writer;
        session.users.forEach((user) => {
            user.own = false
            if (user.id === userId) {
                user.own = true
                userName = user.name
            }
            if (user.id === session.writer) {
                writer = user;
            } else {
                spectators.push(user)
            }
        });
    }
    let html = '';
    if (writer !== undefined) {
        html += ', writer: ';
        if (writer.own) {
            html += '<a href="#" id="own-name" onclick="ownUserNameOnclick()">';
        }
        html += writer.name;
        if (writer.own) {
            html += '</a>';
        }
    }
    if (spectators.length > 0) {
        html += ', spectators: ';
        spectators.forEach((user, i) => {
            if (user.own) {
                html += '<a href="#" id="own-name" onclick="ownUserNameOnclick()">';
            }
            if (i > 0) {
                html += ', ';
            }
            html += user.name;
            if (user.own) {
                html += '</a>';
            }
        })
    }
    if (usersContainerBlock.innerHTML !== html) {
        usersContainerBlock.innerHTML = html;
    }
};
updateUsers();

let writerBlocksUpdate = () => {
    becomeWriterButton.style.display = !isWriter ? 'block' : 'none';
    langSelect.style.display = isWriter ? 'block' : 'none';
    codeBlock.setOption('readOnly', !isWriter);
};
writerBlocksUpdate();

let runnerBlocksUpdate = () => {
    if (session.isRunnerOnline) {
        runnerContainerBlock.style.display = 'none';
    }
    runnerEditButton.style.display = session.isRunnerOnline ? 'none' : 'block';
    runButton.style.display = session.isRunnerOnline && isWriter ? 'block' : 'none';
};
runnerBlocksUpdate();
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
        codeContainerBlock.style.width = null;
        resultContainerBlock.style.display = 'block';
        if (resultBlock.getValue().startsWith('In progress')) {
            resultBlock.setValue(resultBlock.getValue() + '.');
        } else {
            resultBlock.setValue('In progress...');
        }
    } else if (session.result.length > 0) {
        codeContainerBlock.style.width = null;
        resultContainerBlock.style.display = 'block';
        if (sessionPreviousState.result.hash() !== session.result.hash()) {
            resultBlock.setValue(session.result);
        }
    } else if (session.isRunnerOnline) {
        codeContainerBlock.style.width = null;
        resultContainerBlock.style.display = 'block';
        resultBlock.setValue('runner will write result here...');
    } else {
        codeContainerBlock.style.width = '95vw';
        resultContainerBlock.style.display = 'none';
        resultBlock.setValue('...');
    }
};
resultBlockUpdate();

let lastUpdateTimestamp = +new Date;
let pageUpdater = () => {
    let start = +new Date;
    postRequest('/action/session.php', {
        session: session.id,
        user: userId,
        userName: userName,
        lastUpdate: session.updatedAt ? session.updatedAt.data : null,
        action: 'getUpdate',
    }, (response) => {
        ping = +new Date - start;
        console.log((new Date).toLocaleString() + ' | ping: ' + ping);
        lastUpdateTimestamp = +new Date;
        if (response.length === 0) {
            resultBlockUpdate(); // adding more dots to "In progress..."
            return;
        }
        let data = JSON.parse(response);
        if (data.error !== undefined) {
            console.log('getUpdate error', data);
            return
        }

        isNewSession = false;
        sessionPreviousState = {...session};
        session = data;

        // update users
        updateUsers();

        // update writer/spectator ui
        writerBlocksUpdate();

        // update runner ui
        runnerBlocksUpdate();

        // update result ui
        resultBlockUpdate();

        // update session name
        if (sessionPreviousState.name !== session.name) {
            sessionNameBlock.innerHTML = session.name;
        }

        // update code
        if (!isWriter && sessionPreviousState.code.hash() !== session.code.hash()) {
            // if writer, not update code
            let scrollInfo = codeBlock.getScrollInfo();
            codeBlock.setValue(session.code);
            codeBlock.scrollTo(scrollInfo.left, scrollInfo.top);
        }

        // update lang
        if (sessionPreviousState.lang !== session.lang) {
            langSelect.value = session.lang;
            codeBlock.setOption('mode', langKeyToHighlighter[session.lang]);
        }
    }, () => {
        setTimeout(() => {
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
    actions.setCode(() => {
    });
};
pageUpdater();
