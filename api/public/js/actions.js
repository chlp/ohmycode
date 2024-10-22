let actions = {
    setSessionName: () => {
        let newSessionName = sessionNameBlock.textContent;
        postRequest('/action/session.php?action=set_session_name', {
            session: sessionId,
            user: userId,
            userName: userName,
            lang: currentLang,
            action: 'set_session_name',
            sessionName: newSessionName,
        }, (response) => {
            console.log('setSessionName: result', newSessionName, response);
        }, () => {
        });
    },
    setUserName: () => {
        let newUserName = userOwnNameBlock.textContent;
        postRequest('/action/session.php?action=set_user_name', {
            session: sessionId,
            user: userId,
            userName: newUserName,
            action: 'set_user_name',
            lang: currentLang,
        }, (response) => {
            console.log('setUserName: result', newUserName, response);
            if (response === '') {
                localStorage['initialUserName'] = newUserName;
            }
        }, () => {
        });
    },
    setLang: (lang) => {
        postRequest('/action/session.php?action=set_lang', {
            session: sessionId,
            user: userId,
            userName: userName,
            action: 'set_lang',
            lang: lang,
        }, (response) => {
            console.log('setLang: result', response);
            if (response === '') {
                lang = currentLang;
                localStorage['initialLang'] = lang;
            }
        }, () => {
        });
    },
    setRunner: () => {
        postRequest('/action/session.php?action=set_runner', {
            session: sessionId,
            user: userId,
            userName: userName,
            action: 'set_runner',
            runner: runnerInput.value,
            lang: currentLang,
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
        postRequest('/action/session.php?action=set_code', {
            session: sessionId,
            user: userId,
            userName: userName,
            action: 'set_code',
            code: newCode,
            lang: currentLang,
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