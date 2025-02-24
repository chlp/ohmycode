let actions = {
    setFileName: (newFileName) => {
        postRequest('set_name', {
            file_name: newFileName,
        }, () => {
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
    setContent: (content, callback) => {
        if (!isOnline) {
            callback();
            return;
        }
        if (file.writer_id !== '' && file.writer_id !== appId) {
            callback();
            return;
        }
        file.writer_id = appId;
        file.content = content;
        postRequest('set_content', {
            content: content,
        }, callback);
    },
    cleanResult: (callback) => {
        postRequest('clean_result', {}, callback);
    },
    runTask: () => {
        postRequest('run_task', {});
    },
};