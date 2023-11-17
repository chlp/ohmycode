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

let codeBlock = CodeMirror.fromTextArea(document.getElementById("code"), {
    lineNumbers: true,
    mode: langKeyToHighlighter[session.lang], // javascript, go, php, sql
    matchBrackets: true,
    indentWithTabs: false,
    tabSize: 4,
});
let resultBlock = CodeMirror.fromTextArea(document.getElementById("result"), {
    lineNumbers: true,
    readOnly: true,
});

let usersContainerBlock = document.getElementById('users-container');
let sessionNameBlock = document.getElementById('session-name');
let sessionNameInput = document.getElementById('session-name-input');
let userNameContainerBlock = document.getElementById('user-name-container');
let userNameInput = document.getElementById('user-name-input');
let sessionNameContainerBlock = document.getElementById('session-name-container');
let sessionStatusBlock = document.getElementById('session-status');
let becomeWriterButton = document.getElementById('become-writer-button');
let langSelect = document.getElementById('lang-select');
let executeButton = document.getElementById('execute-button');
let executorContainerBlock = document.getElementById('executor-container');
let executorInput = document.getElementById('executor-input');

sessionNameBlock.onclick = () => {
    if (sessionNameContainerBlock.style.display === 'block') {
        sessionNameContainerBlock.style.display = 'none';
    } else {
        sessionNameInput.value = session.name;
        sessionNameContainerBlock.style.display = 'block';
    }
};
let ownUserNameOnclick = () => {
    if (userNameContainerBlock.style.display === 'block') {
        userNameContainerBlock.style.display = 'none';
    } else {
        userNameInput.value = userName;
        userNameContainerBlock.style.display = 'block';
    }
};
langSelect.onchange = () => {
    if (isWriter) {
        codeBlock.setOption('mode', langKeyToHighlighter[langSelect.value]);
        setLang();
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

let executorBlocksUpdate = () => {
    executorContainerBlock.style.display = !session.isExecutorOnline ? 'block' : 'none';
    executeButton.style.display = session.isExecutorOnline && isWriter ? 'block' : 'none';
};
executorBlocksUpdate();

let resultBlockUpdate = () => {
    if (session.isWaitingForResult) {
        resultBlock.setValue('In progress...');
    } else if (session.result.length > 0) {
        if (sessionPreviousState.result.hash() !== session.result.hash()) {
            resultBlock.setValue(session.result);
        }
    } else if (session.isExecutorOnline) {
        resultBlock.setValue('Waiting for execution...');
    } else {
        resultBlock.setValue('No executor');
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
        action: 'getUpdate',
    }, (response) => {
        ping = +new Date - start;
        console.log((new Date).toLocaleString() + ' | ping: ' + ping);
        lastUpdateTimestamp = +new Date;
        if (response.length === 0) {
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

        // update executor ui
        executorBlocksUpdate();

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
