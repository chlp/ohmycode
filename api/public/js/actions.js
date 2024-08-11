let actions = {
    setSessionName: () => {
        let newSessionName = sessionNameBlock.textContent;
        postRequest('/action/session.php?action=setSessionName', {
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
        postRequest('/action/session.php?action=setUserName', {
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
        postRequest('/action/session.php?action=setLang', {
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
        postRequest('/action/session.php?action=setRunner', {
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
    setWriter: (callback) => {
        isSetWriterInProgress = true;
        isWriter = true;
        session.writer = userId;
        postRequest('/action/session.php?action=setWriter', {
            session: session.id,
            user: userId,
            userName: userOwnNameBlock.textContent,
            action: 'setWriter',
        }, (response) => {
            if (response.length !== 0) {
                let data = JSON.parse(response);
                if (data.error !== undefined) {
                    console.log('setWriter: error', data);
                    return;
                }
                session.writer = data.writer;
                isWriter = userId === session.writer;
            }
            isSetWriterInProgress = false;
            writerBlocksUpdate();
            updateUsers();
        }, () => {
            isSetWriterInProgress = false;
            callback();
        });
    },
    setCode: (callback) => {
        if (session.writer !== '' && session.writer !== userId) {
            callback();
            return;
        }
        let sendRequest = () => {
            session.code = codeBlock.getValue();
            postRequest('/action/session.php?action=setCode' , {
                session: session.id,
                user: userId,
                userName: userName,
                action: 'setCode',
                code: codeBlock.getValue(),
            }, (response) => {
                console.log('setCode: result', response);
            }, () => {
                callback();
            });
        }
        if (isWriter) {
            sendRequest();
        } else {
            actions.setWriter(() => {
                if (isWriter) {
                    sendRequest();
                }
            });
        }
    },
    runCode: (callback) => {
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
    },
};