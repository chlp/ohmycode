let dbInstance = null;
const openDB = () => {
    if (dbInstance) {
        return Promise.resolve(dbInstance);
    }

    return new Promise((resolve, reject) => {
        const request = indexedDB.open("FilesDB", 5);

        request.onupgradeneeded = (event) => {
            let db = event.target.result;

            let store;
            if (!db.objectStoreNames.contains("files")) {
                store = db.createObjectStore("files", {keyPath: "id"});
            } else {
                store = event.target.transaction.objectStore("files");
            }

            if (!store.indexNames.contains("content_updated_at")) {
                store.createIndex("content_updated_at", "content_updated_at", {unique: false});
            }
        };

        request.onsuccess = (event) => {
            dbInstance = event.target.result;
            resolve(dbInstance);
        };

        request.onerror = () => reject("Error opening database");
    });
};

const saveFileToDB = async (file) => {
    const db = await openDB();
    const tx = db.transaction("files", "readwrite");
    const store = tx.objectStore("files");
    store.put({
        id: file.id,
        name: file.name,
        lang: file.lang,
        runner: file.runner,
        is_runner_online: file.is_runner_online,
        updated_at: file.updated_at,
        content_updated_at: file.content_updated_at,
        users: file.users,
        is_waiting_for_result: file.is_waiting_for_result,
        result: file.result,
        persisted: file.persisted,
        writer_id: file.writer_id,
        content: file.content,
    });
    return tx.complete;
};

const getFileFromDB = async (id) => {
    const db = await openDB();
    return new Promise((resolve, reject) => {
        const tx = db.transaction("files", "readonly");
        const store = tx.objectStore("files");
        const request = store.get(id);
        request.onsuccess = () => resolve(request.result);
        request.onerror = () => reject("Error fetching file by id: " + id);
    });
};

const getSortedFilesFromDB = async () => {
    const db = await openDB();
    return new Promise((resolve, reject) => {
        const tx = db.transaction("files", "readonly");
        const store = tx.objectStore("files");
        const index = store.index("content_updated_at");
        const request = index.getAll();

        request.onsuccess = () => {
            resolve(request.result.sort((a, b) => new Date(b.content_updated_at) - new Date(a.content_updated_at)));
        };
        request.onerror = () => reject("Error fetching files");
    });
};

const deleteFileInDB = async (id) => {
    const db = await openDB();
    const tx = db.transaction("files", "readwrite");
    tx.objectStore("files").delete(id);
    return tx.complete;
};

export {saveFileToDB, getFileFromDB, getSortedFilesFromDB, deleteFileInDB};