

let socket = null;
let createWebSocket = () => {
    socket = new WebSocket(`${apiUrl}/file`);
    socket.onopen = () => {
        socket.send(JSON.stringify({
            action: 'init',
            file_id: file.id,
            app_id: app.id,
            user_id: app.userId,
            user_name: app.userName,
            lang: app.lang,
        }));
    };
    socket.onclose = (event) => {
        if (event.wasClean) {
            console.log(`Connection closed, code=${event.code}, reason=${event.reason}`);
        } else {
            console.log('Connection closed with error');
        }
        app.isOnline = false;
        writerBlocksUpdate();
        socket = null;
    };
    socket.onerror = (error) => {
        console.log('WebSocket error: ', error);
    };
    socket.onmessage = (event) => {
        try {
            const data = JSON.parse(event.data);
            if (data.error !== undefined) {
                console.log('file::pageUpdater: getUpdate error', data);
                return;
            }

            let previousWriterId = file.writer_id;

            if (typeof data.content === 'undefined') {
                data.content = file.content;
            }
            file = data;

            if (file.persisted) {
                FilesHistory.saveFileToDB(file.id, file.name, file.content_updated_at);
            }
            document.title = `OhMyCode – ${file.name}`;

            // update users
            updateUsers();

            // update "Code is writing now by" block
            writerBlocksUpdate();

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
                writerBlocksUpdate();
            }
        } catch (error) {
            console.error('Wrong message:', error);
        }
    };
};
createWebSocket();
let reconnectAttempts = 0;
setInterval(() => {
    if (socket === null) {
        reconnectAttempts++;
        createWebSocket();
    } else {
        reconnectAttempts = 0;
    }
}, 1000 * Math.min(2 ** reconnectAttempts, 30) + Math.random() * 3000);
