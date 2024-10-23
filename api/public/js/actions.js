let actions = {
    setSessionName: () => {
        let newSessionName = sessionNameBlock.textContent;
        postRequest('/file/set_name', {
            file_id: fileId,
            user_id: userId,
            user_name: userName,
            lang: currentLang,
            session_name: newSessionName,
        }, (response) => {
            console.log('setSessionName: result', newSessionName, response);
        }, () => {
        });
    },
    setUserName: () => {
        let newUserName = userOwnNameBlock.textContent;
        postRequest('/file/set_user_name', {
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
        postRequest('/file/set_lang', {
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
        postRequest('/file/set_runner', {
            file_id: fileId,
            user_id: userId,
            user_name: userName,
            runner_id: runnerInput.value,
            lang: currentLang,
        }, (response) => {
            console.log('setRunner: result', response);
        }, () => {
        });
    },
    setCode: (callback) => {
        if (file.writer_id !== '' && file.writer_id !== userId) {
            callback();
            return;
        }
        file.writer_id = userId;
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
                if (file.writer_id === userId) {
                    file.writer_id = '?';
                }
            }
        }, () => {
            callback();
        });
    },
    cleanResult: (callback) => {
        postRequest('/result/clean', {
            file_id: fileId,
            user_id: userId,
        }, (response, statusCode) => {
            if (statusCode !== 200) {
                console.log('cleanCode: result', response, statusCode);
            }
        }, () => {
            callback();
        });
    },
    runCode: (callback) => {
        postRequest('/run/add_task', {
            file_id: fileId,
            user_id: userId,
        }, (response, statusCode) => {
            if (statusCode !== 200) {
                console.log('runCode: result', response, statusCode);
            }
        }, () => {
            callback();
        });
    },
};