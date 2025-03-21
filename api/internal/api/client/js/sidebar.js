import {openFile} from "./app.js";

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
    getSortedFilesFromDB().then(historyFiles => {
        let htmlLines = historyFiles.map(historyFile =>
            `<span class="history-item">` +
            `<a class="history-delete" data-file-id="${historyFile.id}">x</a> ` +
            `<a class="history-go" data-file-id="${historyFile.id}">${historyFile.name}</a>` +
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

const sidebarBlock = document.getElementById('sidebar');
const historyBlock = document.getElementById('history');
const collapseWithSidebarBlocks = document.getElementsByClassName('collapse-with-sidebar');

let isSidebarVisible = true;
if (localStorage['isSidebarVisible'] === undefined) {
    localStorage['isSidebarVisible'] = JSON.stringify(isSidebarVisible);
} else {
    isSidebarVisible = JSON.parse(localStorage['isSidebarVisible']);
}
const showSidebar = () => {
    sidebarBlock.style.flexBasis = '18em';
    for (const block of collapseWithSidebarBlocks) {
        block.innerHTML = block.dataset.fullText;
    }
    historyBlock.style.display = '';
    setTimeout(() => {
        historyBlock.style.opacity = '1';
    }, 1);
    isSidebarVisible = true;
    localStorage['isSidebarVisible'] = JSON.stringify(true);
};
const hideSidebar = () => {
    sidebarBlock.style.flexBasis = '3em';
    for (const block of collapseWithSidebarBlocks) {
        block.innerHTML = block.dataset.collapsedText;
    }
    historyBlock.style.opacity = '0';
    setTimeout(() => {
        historyBlock.style.display = 'none';
    }, 500); // 500 - the same as the CSS style #sidebar transition: width 0.5s ease;
    isSidebarVisible = false;
    localStorage['isSidebarVisible'] = JSON.stringify(false);
};
const sidebarToggleVisibilitySpan = document.getElementById('sidebar-toggle-visibility');
sidebarToggleVisibilitySpan.onclick = () => {
    if (isSidebarVisible) {
        hideSidebar();
    } else {
        showSidebar();
    }
};
if (isSidebarVisible) {
    showSidebar();
} else {
    hideSidebar();
}

document.addEventListener("DOMContentLoaded", () => {
    document.getElementById("history").addEventListener("click", async (event) => {
        if (event.target.classList.contains("history-delete")) {
            const fileId = event.target.dataset.fileId;
            await deleteFileInDB(fileId);
        } else if (event.target.classList.contains("history-go")) {
            const fileId = event.target.dataset.fileId;
            openFile(fileId, true);
        }
    });
});

export {saveFileToDB};