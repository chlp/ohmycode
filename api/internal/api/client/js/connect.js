import {app, file} from "./app.js";
import {applyFile} from "./file.js";
import {getCurrentLang} from "./lang.js";

let socket = null;
let versionsHandler = null;
let openFileHandler = null;

const onVersions = (handler) => {
    versionsHandler = handler;
};

const onOpenFile = (handler) => {
    openFileHandler = handler;
};

const postRequest = (action, data, callback) => {
    if (socket !== null) {
        socket.send(JSON.stringify({
            ...data,
            action: action,
        }));
    }
    if (typeof callback === 'function') {
        callback();
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
        console.log('openFile');
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
    runTaskWithContent: (content) => {
        if (!app.isOnline) {
            return;
        }
        // Keep local file state in sync to prevent background content sender
        // from firing a separate set_content right before/after Run.
        if (file.writer_id !== '' && file.writer_id !== app.id) {
            return;
        }
        file.writer_id = app.id;
        file.content = content;
        postRequest('run_task_with_content', {
            content: content,
        });
    },
    getVersions: () => {
        postRequest('get_versions', {});
    },
    restoreVersion: (versionId) => {
        postRequest('restore_version', {
            version_id: versionId,
        });
    },
};

const doConnect = (app) => {
    if (socket !== null) {
        actions.openFile();
        return;
    }

    socket = new WebSocket(`${apiUrl}/file`);
    socket.onopen = () => {
        console.log(`Connection opened`);
        actions.openFile();
    };
    socket.onclose = (event) => {
        if (event.wasClean) {
            console.log(`Connection closed, code=${event.code}, reason=${event.reason}`);
        } else {
            console.log('Connection closed with error');
        }
        app.isOnline = false;
        socket = null;
    };
    socket.onerror = (error) => {
        console.log('WebSocket error: ', error);
    };
    socket.onmessage = (event) => {
        try {
            const data = JSON.parse(event.data);
            if (data.error !== undefined) {
                console.log('onmessage: wrong data', data);
                return;
            }

            if (data.action === 'versions') {
                console.log('versions from server');
                if (versionsHandler) {
                    versionsHandler(data.versions);
                }
                return;
            }

            if (data.action === 'open_file') {
                console.log('open_file from server:', data.file_id, 'handler:', !!openFileHandler);
                if (openFileHandler) {
                    openFileHandler(data.file_id);
                } else {
                    console.log('openFileHandler not registered!');
                }
                return;
            }

            console.log('file from server');

            applyFile(data);

            app.isOnline = true;
        } catch (error) {
            console.error('Wrong message:', error);
        }
    };
};

let reconnectAttempts = 0;
setInterval(() => {
    if (socket === null) {
        reconnectAttempts++;
        doConnect(app);
    } else {
        reconnectAttempts = 0;
    }
}, 1000 * Math.min(2 ** reconnectAttempts, 30) + 3000);

export {actions, doConnect, onVersions, onOpenFile};