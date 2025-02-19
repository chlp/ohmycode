const FilesHistory = (() => {
    let dbInstance = null;
    const openDB = () => {
        if (dbInstance) {
            return Promise.resolve(dbInstance); // Если уже есть соединение, используем его
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

    // monitor changes in the db
    const updateHistoryBlock = () => {
        getSortedFilesFromDB().then(files => {
            let htmlLines = files.map(file =>
                `<span class="history-item">` +
                `<a class="history-delete" onclick="FilesHistory.deleteFileInDB('${file.id}')">x</a> ` +
                `<a href="/${file.id}">${file.name}</a>` +
                `</span>`
            );
            historyBlock.innerHTML = htmlLines.join("<br>");
        });
    };
    updateHistoryBlock();

    let lastUpdate = null;
    let historyFilesCount = 0;

    const checkForUpdates = async () => {
        const files = await getSortedFilesFromDB();
        let latestTimestamp = '';
        if (files.length > 0) {
            latestTimestamp = files[0].updated_at;
            historyFilesCount = files.length;
        }
        if (latestTimestamp !== lastUpdate || historyFilesCount !== files.length) {
            lastUpdate = latestTimestamp;
            historyFilesCount = files.length;
            updateHistoryBlock();
        }
    }

    setInterval(checkForUpdates, 1000);

    const deleteFileInDB = async (id) => {
        const db = await openDB();
        const tx = db.transaction("files", "readwrite");
        tx.objectStore("files").delete(id);
        setTimeout(updateHistoryBlock, 100);
        return tx.complete;
    };

    // toggle history button
    const historyBlock = document.getElementById('history');
    const fileBlock = document.getElementById('file');

    let isHistoryVisible = false;
    if (localStorage['isHistoryVisible'] === undefined) {
        localStorage['isHistoryVisible'] = JSON.stringify(isHistoryVisible);
    } else {
        isHistoryVisible = JSON.parse(localStorage['isHistoryVisible']);
    }
    const showHistory = () => {
        historyBlock.style.width = '20em';
        fileBlock.style.width = 'calc(-22em + 100vw)';
        isHistoryVisible = true;
        localStorage['isHistoryVisible'] = JSON.stringify(true);
    };
    const hideHistory = () => {
        historyBlock.style.width = '0';
        fileBlock.style.width = 'calc(-2em + 100vw)';
        isHistoryVisible = false;
        localStorage['isHistoryVisible'] = JSON.stringify(false);
    };
    const toggleHistoryVisibility = () => {
        if (isHistoryVisible) {
            hideHistory();
        } else {
            showHistory();
        }
    };
    if (isHistoryVisible) {
        showHistory();
    }

    return {
        saveFileToDB,
        deleteFileInDB,
        toggleHistoryVisibility
    };
})();