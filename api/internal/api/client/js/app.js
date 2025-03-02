let file = {
    id: getFileIdFromUrl(),
    name: "",
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

let app = {
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
    lang: undefined,
    renderer: undefined,
};

if (localStorage['user_id'] === undefined) {
    localStorage['user_id'] = app.userId;
}
