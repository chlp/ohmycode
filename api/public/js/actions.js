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
    setCode: (callback) => {
        if (session.writer !== '' && session.writer !== userId) {
            callback();
            return;
        }
        session.writer = userId;
        let newCode = codeBlock.getValue();
        session.code = newCode;
        postRequest('/action/session.php?action=setCode', {
            session: session.id,
            user: userId,
            userName: userName,
            action: 'setCode',
            code: newCode,
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
    runCode: (callback) => {
        session.result = 'In progress..';
        resultBlock.setValue('In progress..');
        postRequest('/action/request.php?action=set', {
            session: session.id,
            action: 'set',
        }, (response) => {
            console.log('runCode->setCode: result', response);
        }, () => {
            callback();
        });
    },
};