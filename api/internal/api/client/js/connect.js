import {app, file} from "./app.js";
import {applyFile} from "./file.js";
import {getCurrentLang} from "./lang.js";
import {encryptText, decryptText} from "./encrypt.js";

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
            ro_token: app.roToken || '',
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
    setContent: async (content) => {
        if (!app.isOnline) {
            return;
        }
        if (file.is_locked) {
            return;
        }
        if (file.writer_id !== '' && file.writer_id !== app.id) {
            return;
        }
        file.writer_id = app.id;
        file.content = content;
        let payload = content;
        let roPayload = '';
        if (file.encrypted) {
            if (app.encKey) {
                try {
                    payload = await encryptText(app.encKey, content);
                } catch(e) {
                    console.error('Encryption failed, not sending:', e);
                    return;
                }
            }
            if (app.roEncKey) {
                try {
                    roPayload = await encryptText(app.roEncKey, content);
                } catch(e) {
                    console.error('RO encryption failed:', e);
                }
            }
        }
        postRequest('set_content', {content: payload, ro_content: roPayload});
    },
    setLocked: (isLocked) => {
        postRequest('set_locked', {
            is_locked: isLocked,
        });
    },
    setEncrypted: (encrypted) => {
        postRequest('set_encrypted', {encrypted});
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
    socket.onmessage = async (event) => {
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

            if (data.encrypted) {
                if (app.isROLink && app.encKey) {
                    if (typeof data.ro_content === 'string' && data.ro_content) {
                        try {
                            data.content = await decryptText(app.encKey, data.ro_content);
                        } catch(e) {
                            console.error('RO decryption failed:', e);
                            delete data.content;
                        }
                    } else {
                        delete data.content;
                    }
                } else if (!app.isROLink && app.encKey && typeof data.content === 'string') {
                    try {
                        data.content = await decryptText(app.encKey, data.content);
                    } catch(e) {
                        console.error('Decryption failed:', e);
                        delete data.content;
                    }
                } else if (typeof data.content !== 'undefined') {
                    delete data.content;
                }
            }

            applyFile(data);

            app.isOnline = true;
        } catch (error) {
            console.error('Wrong message:', error);
        }
    };
};

let reconnectAttempts = 0;
const scheduleReconnect = () => {
    const delayMs = reconnectAttempts === 0
        ? 1000
        : Math.min(1000 * 2 ** reconnectAttempts, 30000);
    setTimeout(() => {
        if (socket === null) {
            reconnectAttempts++;
            doConnect(app);
        } else {
            reconnectAttempts = 0;
        }
        scheduleReconnect();
    }, delayMs);
};
scheduleReconnect();

export {actions, doConnect, onVersions, onOpenFile};