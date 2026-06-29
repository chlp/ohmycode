import {contentCodeMirror, contentMarkdownBlock, updateEditorLockStatus} from "./editor.js";
import {setLang} from "./lang.js";
import {getFileFromDB} from "./db.js";
import {applyFile} from "./file.js";
import {fileNameBlock, fileNameEditing} from "./file_name.js";
import {doConnect} from "./connect.js";
import {importKey} from "./encrypt.js";

const genUuid = () => {
    const alphabet = '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz';
    const bytes = new Uint8Array(22);
    crypto.getRandomValues(bytes);
    return Array.from(bytes, b => alphabet[b % 62]).join('');
};

const isUuid = (id) => {
    return /^[0-9A-Za-z]{22}$/.test(id) || /^[a-z0-9]{32}$/.test(id);
};

const getFileIdFromWindowLocation = () => {
    const url = new URL(window.location.href);
    const fileId = url.pathname.replace(/^\/|\/$/g, '');
    if (!isUuid(fileId)) {
        return undefined;
    }
    return fileId;
};

// isRO: true when opening via a read-only link — determines which localStorage slot to use.
const loadKeyForFile = async (fileId, isRO) => {
    const storagePrefix = isRO ? 'ohmycode_rokey_' : 'ohmycode_key_';

    const currentFileId = getFileIdFromWindowLocation();
    if (currentFileId === fileId && window.location.hash) {
        const hashStr = window.location.hash.slice(1);
        const params = new URLSearchParams(hashStr);
        const keyStr = params.get('key') || (hashStr.length >= 40 ? hashStr : null);
        if (keyStr) {
            try {
                const key = await importKey(keyStr);
                localStorage.setItem(storagePrefix + fileId, keyStr);
                history.replaceState(null, '', window.location.pathname);
                return key;
            } catch(e) { console.warn('Key from URL hash invalid:', e); }
        }
    }
    const stored = localStorage.getItem(storagePrefix + fileId);
    if (stored) {
        try { return await importKey(stored); } catch(e) { console.warn('Stored key invalid, clearing:', e); localStorage.removeItem(storagePrefix + fileId); }
    }
    return null;
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
        encrypted: false,
        ro_token: '',

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

        _is_locked: false,
        get is_locked() {
            return this._is_locked;
        },
        set is_locked(value) {
            if (this._is_locked !== value) {
                this._is_locked = value;
                updateEditorLockStatus();
            }
        },
    };
};

let file;

// Detect read-only access for a file from localStorage.
// Edit key takes priority: if the user has the edit key, they get full access
// even if a ro_token was previously stored (e.g. owner testing their own RO link).
const detectROFromStorage = (id) => {
    const hasEditKey = !!localStorage.getItem('ohmycode_key_' + id);
    const storedRoToken = localStorage.getItem('ohmycode_ro_token_' + id);
    if (!hasEditKey && storedRoToken) {
        return storedRoToken;
    }
    return null;
};

const openFile = async (id, pushHistory) => {
    app.isOnline = false;

    if (pushHistory) {
        app.roToken = null;
        app.isROLink = false;
        localStorage.removeItem('ohmycode_ro_token_' + id);
    } else {
        const storedRoToken = detectROFromStorage(id);
        app.roToken = storedRoToken;
        app.isROLink = !!storedRoToken;
    }

    app.encKey = await loadKeyForFile(id, app.isROLink);

    // Load readonly encryption key for non-RO owners who have generated readonly links.
    app.roEncKey = null;
    if (!app.isROLink) {
        const roKeyStr = localStorage.getItem('ohmycode_rokey_' + id);
        if (roKeyStr) {
            try { app.roEncKey = await importKey(roKeyStr); } catch(e) { localStorage.removeItem('ohmycode_rokey_' + id); }
        }
    }

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
    setLang(localStorage['initialLang']);

    let fileId = getFileIdFromWindowLocation();
    let pushHistory = false;
    if (fileId === undefined) {
        fileId = genUuid();
        pushHistory = true;
    }

    // On first visit via a RO link the ro_token is in the hash (#key=...&ro=TOKEN).
    // Persist it to localStorage so subsequent plain visits stay in read-only mode.
    if (!pushHistory) {
        const hashStr = window.location.hash.slice(1);
        const hashParams = new URLSearchParams(hashStr);
        const roTokenFromHash = hashParams.get('ro');
        if (roTokenFromHash) {
            localStorage.setItem('ohmycode_ro_token_' + fileId, roTokenFromHash);
        }
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
    encKey: null,
    roEncKey: null,
    isROLink: false,
    roToken: null,
};

if (localStorage['user_id'] === undefined) {
    localStorage['user_id'] = app.userId;
}

export {file, app, openFile, loadKeyForFile};
