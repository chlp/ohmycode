let actions = {
    setSessionName: () => {
        let newSessionName = sessionNameBlock.textContent;
        postRequest('/action/session.php?action=set_session_name', {
            file_id: fileId,
            user_id: userId,
            user_name: userName,
            lang: currentLang,
            sessionName: newSessionName,
        }, (response) => {
            console.log('setSessionName: result', newSessionName, response);
        }, () => {
        });
    },
    setUserName: () => {
        let newUserName = userOwnNameBlock.textContent;
        postRequest('/action/session.php?action=set_user_name', {
            file_id: fileId,
            user_id: userId,
            user_name: newUserName,
            lang: currentLang,
        }, (response) => {
            console.log('setuser_name: result', newUserName, response);
            if (response === '') {
                localStorage['initialUserName'] = newUserName;
            }
        }, () => {
        });
    },
    setLang: (lang) => {
        postRequest('/action/session.php?action=set_lang', {
            file_id: fileId,
            user_id: userId,
            user_name: userName,
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
            file_id: fileId,
            user_id: userId,
            user_name: userName,
            runner: runnerInput.value,
            lang: currentLang,
        }, (response) => {
            console.log('setRunner: result', response);
        }, () => {
        });
    },
    setCode: (callback) => {
        if (file.writer !== '' && file.writer !== userId) {
            callback();
            return;
        }
        file.writer = userId;
        let newContent = contentBlock.getValue();
        file.content = newContent;
        postRequest('/file/set_content', {
            file_id: fileId,
            user_id: userId,
            user_name: userName,
            content: newContent,
            lang: currentLang,
        }, (response, statusCode) => {
            if (statusCode !== 200) {
                console.log('setCode: result', response, statusCode);
            }
            if (statusCode === 403) {
                if (file.writer === userId) {
                    file.writer = '?';
                }
            }
        }, () => {
            callback();
        });
    },
    cleanCode: (callback) => {
        postRequest('/action/result.php?action=clean', {
            file_id: fileId,
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
            file_id: fileId,
        }, (response, statusCode) => {
            if (statusCode !== 200) {
                console.log('runCode: result', response, statusCode);
            }
        }, () => {
            callback();
        });
    },
};