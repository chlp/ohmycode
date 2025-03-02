let file = {
    "id": getFileIdFromUrl(),
    "name": "",
    "content": "",
    "lang": 'markdown',
    "runner": "",
    "is_runner_online": false,
    "updated_at": null,
    "content_updated_at": null,
    "writer_id": "",
    "users": [],
    "is_waiting_for_result": false,
    "result": "",
    "persisted": false,
};

let app = {
    isOnline: false,
    id: genUuid(),
    userId: localStorage['user_id'] === undefined ? genUuid() : localStorage['user_id'],
    userName: localStorage['user_name'] === undefined ? '' : localStorage['user_name'],
    lang: undefined,
    renderer: undefined,
};

if (localStorage['user_id'] === undefined) {
    localStorage['user_id'] = app.userId;
}
