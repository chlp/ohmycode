import {ohMySimpleHash} from "./utils.js";
import {app, file, openFile} from "./app.js";
import {saveFileToDB} from "./sidebar.js";
import {contentCodeMirror, contentMarkdownBlock} from "./editor.js";
import {getCurrentLang} from "./lang.js";

const postRequest = (action, data, callback) => {
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

const getLocalDateTimeString = () => {
    return new Intl.DateTimeFormat('sv-SE', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
        hour12: false
    }).format(new Date()).replace(',', '');
};

const actions = {
    openFile: () => {
        postRequest('init', {
            file_id: file.id,
            file_name: getLocalDateTimeString(),
            app_id: app.id,
            user_id: app.userId,
            user_name: app.userName,
            lang: getCurrentLang().key,
        });
    },
    setFileName: (newFileName) => {
        postRequest('set_name', {
            file_name: newFileName,
        });
    },
    setUserName: (newUserName) => {
        postRequest('set_user_name', {
            user_name: newUserName,
        });
    },
    setLang: (lang) => {
        if (!app.isOnline) {
            return;
        }
        postRequest('set_lang', {
            lang: lang,
        });
    },
    setRunner: (runnerId) => {
        postRequest('set_runner', {
            runner_id: runnerId,
        });
    },
    setContent: (content) => {
        if (!app.isOnline) {
            return;
        }
        if (file.writer_id !== '' && file.writer_id !== app.id) {
            return;
        }
        file.writer_id = app.id;
        file.content = content;
        postRequest('set_content', {
            content: content,
        });
    },
    cleanResult: () => {
        postRequest('clean_result', {});
    },
    runTask: () => {
        postRequest('run_task', {});
    },
};

const fileChangeHandlers = [];
const onFileChange = (callback) => {
    if (typeof callback === "function") {
        fileChangeHandlers.push(callback);
    }
};

const controlsContainerBlock = document.getElementById('controls-container');
const createWebSocket = (app) => {
    app.socket = new WebSocket(`${apiUrl}/file`);
    app.socket.onopen = () => {
        console.log(`Connection opened`);
        actions.openFile();
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
                console.log('onmessage: wrong data', data);
                return;
            }

            if (file.id !== data.id) {
                openFile(file.id);
                console.log('onmessage: new file.id', data.id, file.id);
                return;
            }

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
            file.writer_id = data.writer_id;
            if (typeof data.content === 'string') {
                file.content = data.content;
            }

            if (file.persisted) {
                saveFileToDB(file.id, file.name, file.content_updated_at);
            }
            document.title = `OhMyCode – ${file.name}`;

            fileChangeHandlers.forEach(fn => fn(file));

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

            if (controlsContainerBlock.style.display !== 'block') {
                controlsContainerBlock.style.display = 'block';
            }

            if (!app.isOnline) {
                app.isOnline = true;
            }
        } catch (error) {
            console.error('Wrong message:', error);
        }
    };
};
window.addEventListener("DOMContentLoaded", () => {
    setTimeout(()=> {
        createWebSocket(app);
    }, 100);
});

let reconnectAttempts = 0;
setInterval(() => {
    if (app.socket === null) {
        reconnectAttempts++;
        createWebSocket(app);
    } else {
        reconnectAttempts = 0;
    }
}, 1000 * Math.min(2 ** reconnectAttempts, 30) + Math.random() * 3000);

export {actions, onFileChange};