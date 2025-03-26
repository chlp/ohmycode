let dbInstance = null;
const openDB = () => {
    if (dbInstance) {
        return Promise.resolve(dbInstance);
    }

    return new Promise((resolve, reject) => {
        const request = indexedDB.open("FilesDB", 1);

        request.onupgradeneeded = (event) => {
            let db = event.target.result;
            if (!db.objectStoreNames.contains("files")) {
                let store = db.createObjectStore("files", {keyPath: "id"});
                store.createIndex("updated_at", "updated_at", {unique: false});
            }
        };

        request.onsuccess = (event) => {
            dbInstance = event.target.result;
            resolve(dbInstance);
        };
        request.onerror = () => reject("Error opening database");
    });
};

const saveFileToDB = async (id, fileName, updatedAt) => {
    const db = await openDB();
    const tx = db.transaction("files", "readwrite");
    const store = tx.objectStore("files");
    store.put({id: id, name: fileName, updated_at: updatedAt});
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
        const index = store.index("updated_at");
        const request = index.getAll();

        request.onsuccess = () => {
            resolve(request.result.sort((a, b) => new Date(b.updated_at) - new Date(a.updated_at)));
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