let actions = {
    setSessionName: () => {
        let newSessionName = sessionNameBlock.textContent;
        postRequest('/action/session.php', {
            session: session.id,
            user: userId,
            userName: userName,
            action: 'setSessionName',
            sessionName: newSessionName,
        }, (response) => {
            console.log('setSessionName: result', newSessionName, response);
        }, () => {
        });
    },
    setUserName: () => {
        let newUserName = userOwnNameBlock.textContent;
        postRequest('/action/session.php', {
            session: session.id,
            user: userId,
            userName: newUserName,
            action: 'setUserName',
        }, (response) => {
            console.log('setUserName: result', newUserName, response);
            if (response === '') {
                localStorage['initialUserName'] = newUserName;
            }
        }, () => {
        });
    },
    setLang: () => {
        postRequest('/action/session.php', {
            session: session.id,
            user: userId,
            userName: userOwnNameBlock.textContent,
            action: 'setLang',
            lang: langSelect.value,
        }, (response) => {
            console.log('setLang: result', response);
        }, () => {
        });
    },
    setRunner: () => {
        postRequest('/action/session.php', {
            session: session.id,
            user: userId,
            userName: userOwnNameBlock.textContent,
            action: 'setRunner',
            runner: runnerInput.value,
        }, (response) => {
            console.log('setRunner: result', response);
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
            userName: userOwnNameBlock.textContent,
            action: 'setWriter',
        }, (response) => {
            console.log('setWriter: result', response);
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
                console.log('setCode: result', response);
                callback();
            }, () => {
            });
        } else {
            callback();
        }
    },
    runCode: (callback) => {
        actions.setCode(() => {
            session.result = 'In progress..';
            resultBlock.setValue('In progress..');
            postRequest('/action/request.php', {
                session: session.id,
                action: 'set',
            }, (response) => {
                console.log('runCode->setCode: result', response);
            }, () => {
                callback();
            });
        });
    },
};