let actions = {
    setFileName: () => {
        let newFileName = fileNameBlock.textContent;
        postRequest('set_name', {
            file_name: newFileName,
        });
    },
    setUserName: () => {
        let newUserName = userOwnNameBlock.textContent;
        postRequest('set_user_name', {
            user_name: newUserName,
        }, () => {
            localStorage['user_name'] = newUserName;
        });
    },
    setLang: (lang) => {
        postRequest('set_lang', {
            lang: lang,
        }, () => {
            localStorage['initialLang'] = lang;
        });
    },
    setRunner: () => {
        postRequest('set_runner', {
            runner_id: runnerInput.value,
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
        postRequest('set_content', {
            content: newContent,
        }, callback);
    },
    cleanResult: (callback) => {
        postRequest('clean_result', {}, callback);
    },
    runTask: (callback) => {
        postRequest('run_task', {}, callback);
    },
};