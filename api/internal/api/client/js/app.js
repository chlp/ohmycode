import {contentCodeMirror, contentMarkdownBlock, updateEditorLockStatus} from "./editor.js";
import {setLang} from "./lang.js";
import {fileNameBlock, fileNameEditing} from "./file_name.js";
import {actions} from "./connect.js";

const genUuid = () => { // Генерация случайного UUID без дефисов
    return 'xxxxxxxxxxxx4xxxyxxxxxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
        const r = Math.random() * 16 | 0, v = c === 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
};

const isUuid = (id) => {
    return (new RegExp(`^[a-z0-9]{32}$`)).test(id);
};

const getFileIdFromUrl = () => {
    const url = new URL(window.location.href);
    let fileId = url.pathname.replace(/^\/|\/$/g, '');
    if (!isUuid(fileId)) {
        fileId = genUuid();
        history.pushState({}, null, '/' + fileId);
    }
    return fileId;
}

const initFile = () => {
    return {
        id: getFileIdFromUrl(),
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
            if (fileNameBlock.innerHTML !== value && !fileNameEditing) {
                fileNameBlock.innerHTML = value;
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
let file = initFile();

const openFile = (id) => {
    app.isOnline = false;
    file.id = id;
    file.content = "";
    contentCodeMirror.setValue("");
    contentMarkdownBlock.innerHTML = "";
    history.pushState({}, null, '/' + file.id);
    file = initFile();
    actions.openFile();
};

document.getElementById('sidebar-create-new-file').onclick = () => {
    openFile(genUuid());
};

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