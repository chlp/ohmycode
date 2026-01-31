import {contentCodeMirror, contentMarkdownBlock, updateEditorLockStatus} from "./editor.js";
import {setLang} from "./lang.js";
import {getFileFromDB} from "./db.js";
import {applyFile} from "./file.js";
import {fileNameBlock, fileNameEditing} from "./file_name.js";
import {doConnect} from "./connect.js";

const genUuid = () => { // Генерация случайного UUID без дефисов
    return 'xxxxxxxxxxxx4xxxyxxxxxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
        const r = Math.random() * 16 | 0, v = c === 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
};

const isUuid = (id) => {
    return (new RegExp(`^[a-z0-9]{32}$`)).test(id);
};

const getFileIdFromWindowLocation = () => {
    const url = new URL(window.location.href);
    const fileId = url.pathname.replace(/^\/|\/$/g, '');
    if (!isUuid(fileId)) {
        return undefined;
    }
    return fileId;
};

const initFile = (fileId) => {
    if (!isUuid(fileId)) {
        console.error('initFile. Wrong fileId', fileId);
        fileId = genUuid();
    }

    return {
        id: fileId,
        content: "",
        lang: 'markdown',
        runner: "",
        is_runner_online: false,
        updated_at: null,
        content_updated_at: null,
        users: [],
        is_waiting_for_result: false,
        result: "",
        persisted: false,

        _name: "",
        get name() {
            return this._name;
        },
        set name(value) {
            if (fileNameBlock.textContent !== value && !fileNameEditing) {
                fileNameBlock.textContent = value;
            }
            this._name = value;
        },

        _writer_id: "",
        get writer_id() {
            return this._writer_id;
        },
        set writer_id(value) {
            if (this._writer_id !== value) {
                this._writer_id = value;
                updateEditorLockStatus();
            }
        },
    };
};

let file;

const openFile = async (id, pushHistory) => {
    app.isOnline = false;

    file = initFile(id);
    contentCodeMirror.setValue('');
    contentMarkdownBlock.innerHTML = '';
    if (pushHistory) {
        history.pushState({}, null, '/' + file.id);
    }

    let fileFromDb = await getFileFromDB(file.id);
    if (typeof fileFromDb !== 'undefined') {
        console.log('file from db');
        applyFile(fileFromDb); // load really fast
    }

    doConnect(app); // could load longer
};

window.addEventListener("DOMContentLoaded", () => {
    let fileId = getFileIdFromWindowLocation();
    let pushHistory = false;
    if (fileId === undefined) {
        fileId = genUuid();
        pushHistory = true;
    }
    openFile(fileId, pushHistory).then(() => {
    });
});

document.getElementById('header-new-file-btn').onclick = () => {
    openFile(genUuid(), true).then(() => {
    });
};

window.addEventListener("popstate", () => {
    const fileId = getFileIdFromWindowLocation();
    if (fileId !== undefined) {
        openFile(fileId, false).then(() => {
        });
    }
});

const app = {
    _isOnline: false,
    get isOnline() {
        return this._isOnline;
    },
    set isOnline(value) {
        if (this._isOnline !== value) {
            this._isOnline = value;
            updateEditorLockStatus();
        }
    },
    id: genUuid(),
    userId: localStorage['user_id'] === undefined ? genUuid() : localStorage['user_id'],
    userName: localStorage['user_name'] === undefined ? '' : localStorage['user_name'],
    renderer: undefined,
};

if (localStorage['user_id'] === undefined) {
    localStorage['user_id'] = app.userId;
}

setLang(localStorage['initialLang']);

export {file, app, openFile};