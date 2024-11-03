let actions = {
    setFileName: () => {
        let newFileName = fileNameBlock.textContent;
        postRequest('/file/set_name', {
            file_id: fileId,
            user_id: userId,
            user_name: userName,
            lang: currentLang,
            file_name: newFileName,
        }, (response) => {
            console.log('setFileName: result', newFileName, response);
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
        }, (response, statusCode) => {
            console.log('setUserName: result', newUserName, response, statusCode);
            if (statusCode === 200 || statusCode === 204) {
                localStorage['user_name'] = newUserName;
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
    setContent: (callback) => {
        if (!isOnline) {
            callback();
            return;
        }
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
    runTask: (callback) => {
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