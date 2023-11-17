let actions = {
    setSessionName: () => {
        postRequest('/action/session.php', {
            session: session.id,
            user: userId,
            userName: userName,
            action: 'setSessionName',
            sessionName: sessionNameInput.value,
        }, (response) => {
            console.log('saved session name', response);
            sessionNameBlock.innerHTML = sessionNameInput.value;
            sessionNameContainerBlock.style.display = 'none';
        }, () => {
        });
    },
    setUserName: () => {
        postRequest('/action/session.php', {
            session: session.id,
            user: userId,
            userName: userNameInput.value,
            action: 'setUserName',
        }, (response) => {
            console.log('saved user name', response);
            userName = userNameInput.value;
            document.getElementById('own-name').innerHTML = userName;
            userNameContainerBlock.style.display = 'none';
        }, () => {
        });
    },
    setLang: () => {
        postRequest('/action/session.php', {
            session: session.id,
            user: userId,
            userName: userNameInput.value,
            action: 'setLang',
            lang: langSelect.value,
        }, (response) => {
            console.log('saved lang', response);
        }, () => {
        });
    },
    setExecutor: () => {
        postRequest('/action/session.php', {
            session: session.id,
            user: userId,
            userName: userNameInput.value,
            action: 'setExecutor',
            executor: executorInput.value,
        }, (response) => {
            console.log('saved executor', response);
        }, () => {
        });
    },
    setWriter: () => {
        isWriter = true;
        session.writer = userId;
        writerBlocksUpdate();
        updateUsers();
        postRequest('/action/session.php', {
            session: session.id,
            user: userId,
            userName: userNameInput.value,
            action: 'setWriter',
        }, (response) => {
            console.log('saved writer', response);
        }, () => {
        });
    },
    setCode: (callback) => {
        if (!isWriter) {
            return;
        }
        if (session.code.hash() !== codeBlock.getValue().hash()) {
            session.code = codeBlock.getValue();
            postRequest('/action/session.php', {
                session: session.id,
                user: userId,
                userName: userName,
                action: 'setCode',
                code: codeBlock.getValue(),
            }, (response) => {
                console.log('saved code', response);
                callback();
            }, () => {
            });
        } else {
            callback();
        }
    },
    setRequest: () => {
        if (!isWriter) {
            return;
        }
        actions.setCode(() => {
            postRequest('/action/request.php', {
                session: session.id,
                action: 'set',
            }, (response) => {
                console.log('saved request', response);
            }, () => {
            });
        });
    },
};