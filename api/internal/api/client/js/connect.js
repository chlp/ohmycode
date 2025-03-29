import {app, file, openFile} from "./app.js";
import {loadNewFileVersion} from "./file.js";
import {getCurrentLang} from "./lang.js";

const postRequest = (action, data, callback) => {
    try {
        if (typeof app.socket === 'undefined') {
            if (typeof callback === 'function') {
                callback();
            }
        } else {
            app.socket.send(JSON.stringify({
                ...data,
                action: action,
            }));
        }
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
                openFile(file.id, false);
                console.log('onmessage: new file.id', data.id, file.id);
                return;
            }

            loadNewFileVersion(data);

            app.isOnline = true;
        } catch (error) {
            console.error('Wrong message:', error);
        }
    };
};
window.addEventListener("DOMContentLoaded", () => {
    setTimeout(() => {
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
}, 1000 * Math.min(2 ** reconnectAttempts, 30) + 3000);

export {actions};