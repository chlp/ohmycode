let openDB = () => {
    return new Promise((resolve, reject) => {
        const request = indexedDB.open("FilesDB", 1);

        request.onupgradeneeded = (event) => {
            let db = event.target.result;
            if (!db.objectStoreNames.contains("files")) {
                let store = db.createObjectStore("files", {keyPath: "id"});
                store.createIndex("updated_at", "updated_at", {unique: false});
            }
        };

        request.onsuccess = () => resolve(request.result);
        request.onerror = () => reject("Error opening database");
    });
};

let saveFileToDB = async (id, fileName, updatedAt) => {
    const db = await openDB();
    const tx = db.transaction("files", "readwrite");
    const store = tx.objectStore("files");
    store.put({id: id, name: fileName, updated_at: updatedAt});
    return tx.complete;
};

let getSortedFilesFromDB = async () => {
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

// monitor changes in the db
let updateHistoryBlock = () => {
    getSortedFilesFromDB().then(files => {
        let htmlLines = files.map(file => `<a href="/${file.id}">${file.name}</a>`);
        historyBlock.innerHTML = htmlLines.join("<br>");
    });
};
updateHistoryBlock();

let lastUpdate = null;

async function checkForUpdates() {
    const files = await getSortedFilesFromDB();
    let latestTimestamp = '';
    if (files.length > 0) {
        latestTimestamp = files[0].updated_at;
    }
    if (latestTimestamp !== lastUpdate) {
        lastUpdate = latestTimestamp;
        updateHistoryBlock();
    }
}

setInterval(checkForUpdates, 1000);

let deleteFileInDB = async (id) => {
    const db = await openDB();
    const tx = db.transaction("files", "readwrite");
    tx.objectStore("files").delete(id);
    return tx.complete;
};
