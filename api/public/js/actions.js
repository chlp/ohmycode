let actions = {
    setSessionName: () => {
        let newSessionName = sessionNameBlock.textContent;
        postRequest('/action/session.php?action=setSessionName', {
            session: sessionId,
            user: userId,
            userName: userName,
            lang: langSelect.value,
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
            session: sessionId,
            user: userId,
            userName: newUserName,
            action: 'setUserName',
            lang: langSelect.value,
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
            session: sessionId,
            user: userId,
            userName: userName,
            action: 'setLang',
            lang: langSelect.value,
        }, (response) => {
            console.log('setLang: result', response);
            if (response === '') {
                localStorage['initialLang'] = langSelect.value;
            }
        }, () => {
        });
    },
    setRunner: () => {
        postRequest('/action/session.php?action=setRunner', {
            session: sessionId,
            user: userId,
            userName: userName,
            action: 'setRunner',
            runner: runnerInput.value,
            lang: langSelect.value,
        }, (response) => {
            console.log('setRunner: result', response);
        }, () => {
        });
    },
    setCode: (callback) => {
        if (session.writer !== '' && session.writer !== userId) {
            callback();
            return;
        }
        session.writer = userId;
        let newCode = codeBlock.getValue();
        session.code = newCode;
        postRequest('/action/session.php?action=setCode', {
            session: sessionId,
            user: userId,
            userName: userName,
            action: 'setCode',
            code: newCode,
            lang: langSelect.value,
        }, (response, statusCode) => {
            if (statusCode !== 200) {
                console.log('setCode: result', response, statusCode);
            }
            if (statusCode === 403) {
                if (session.writer === userId) {
                    session.writer = '?';
                }
            }
        }, () => {
            callback();
        });
    },
    cleanCode: (callback) => {
        postRequest('/action/result.php?action=clean', {
            session: sessionId,
            action: 'clean',
        }, (response, statusCode) => {
            if (statusCode !== 200) {
                console.log('cleanCode: result', response, statusCode);
            }
        }, () => {
            callback();
        });
    },
    runCode: (callback) => {
        postRequest('/action/request.php?action=set', {
            session: sessionId,
            action: 'set',
        }, (response, statusCode) => {
            if (statusCode !== 200) {
                console.log('runCode: result', response, statusCode);
            }
        }, () => {
            callback();
        });
    },
};