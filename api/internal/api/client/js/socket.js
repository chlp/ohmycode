const controlsContainerBlock = document.getElementById('controls-container');
let createWebSocket = (app) => {
    app.socket = new WebSocket(`${apiUrl}/file`);
    app.socket.onopen = () => {
        console.log(`Connection opened`);
        app.socket.send(JSON.stringify({
            action: 'init',
            file_id: file.id,
            app_id: app.id,
            user_id: app.userId,
            user_name: app.userName,
            lang: app.lang,
        }));
    };
    app.socket.onclose = (event) => {
        if (event.wasClean) {
            console.log(`Connection closed, code=${event.code}, reason=${event.reason}`);
        } else {
            console.log('Connection closed with error');
        }
        app.isOnline = false;
        app.socket = null;
    };
    app.socket.onerror = (error) => {
        console.log('WebSocket error: ', error);
    };
    app.socket.onmessage = (event) => {
        try {
            const data = JSON.parse(event.data);
            if (data.error !== undefined) {
                console.log('file::pageUpdater: getUpdate error', data);
                return;
            }

            // todo: prepare to change file.id and load brand new file

            let previousWriterId = file.writer_id;

            file.name = data.name;
            file.lang = data.lang;
            file.runner = data.runner;
            file.is_runner_online = data.is_runner_online;
            file.updated_at = data.updated_at;
            file.content_updated_at = data.content_updated_at;
            file.users = data.users;
            file.is_waiting_for_result = data.is_waiting_for_result;
            file.result = data.result;
            file.persisted = data.persisted;
            file._writer_id = data.writer_id;
            if (typeof data.content === 'string') {
                file.content = data.content;
            }

            if (file.persisted) {
                FilesHistory.saveFileToDB(file.id, file.name, file.content_updated_at);
            }
            document.title = `OhMyCode – ${file.name}`;

            // update users
            updateUsers();

            // update runner ui
            runnerBlocksUpdate();

            // update result ui
            resultBlockUpdate();

            // update file name
            if (fileNameBlock.innerHTML !== file.name && !fileNameEditing) {
                fileNameBlock.innerHTML = file.name;
            }

            // update code
            if (
                !app.isOnline || // first load
                (
                    file.writer_id !== app.id && previousWriterId !== app.id && // do not update if current user is writer
                    ohMySimpleHash(file.content) !== ohMySimpleHash(contentCodeMirror.getValue()) // do not update if code is the same already
                )
            ) {
                let {left, top} = contentCodeMirror.getScrollInfo();
                let {line, ch} = contentCodeMirror.getCursor();
                contentCodeMirror.setValue(file.content);
                contentMarkdownBlock.innerHTML = marked.parse(file.content);
                contentCodeMirror.scrollTo(left, top);
                contentCodeMirror.setCursor({line: line, ch: ch});
            }

            // update lang
            setLang(file.lang);

            if (controlsContainerBlock.style.display !== 'block') {
                controlsContainerBlock.style.display = 'block';
            }

            if (!app.isOnline) {
                app.isOnline = true;
                contentSender();
            }
        } catch (error) {
            console.error('Wrong message:', error);
        }
    };
};
createWebSocket(app);
let reconnectAttempts = 0;
setInterval(() => {
    if (app.socket === null) {
        reconnectAttempts++;
        createWebSocket(app);
    } else {
        reconnectAttempts = 0;
    }
}, 1000 * Math.min(2 ** reconnectAttempts, 30) + Math.random() * 3000);

let postRequest = (action, data, callback) => {
    try {
        app.socket.send(JSON.stringify({
            ...data,
            action: action,
        }));
    } finally {
        if (typeof callback === 'function') {
            callback();
        }
    }
};