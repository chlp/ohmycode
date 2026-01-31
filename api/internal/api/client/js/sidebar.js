import {openFile} from "./app.js";
import {deleteFileInDB, getSortedFilesFromDB} from "./db.js";

const historyPanel = document.getElementById('history-panel');
const historyBlock = document.getElementById('history');
const historyBtn = document.getElementById('header-history-btn');

// Update history list
const updateHistoryBlock = () => {
    getSortedFilesFromDB().then(historyFiles => {
        historyBlock.innerHTML = '';
        historyFiles.forEach((historyFile) => {
            let span = document.createElement('span');
            span.className = 'history-item';

            let deleteLink = document.createElement('a');
            deleteLink.className = 'history-delete';
            deleteLink.dataset.fileId = historyFile.id;
            deleteLink.textContent = 'x';

            let goLink = document.createElement('a');
            goLink.className = 'history-go';
            goLink.dataset.fileId = historyFile.id;
            goLink.textContent = historyFile.name;

            span.appendChild(deleteLink);
            span.appendChild(document.createTextNode(' '));
            span.appendChild(goLink);
            historyBlock.appendChild(span);
        });
    });
};

window.addEventListener("DOMContentLoaded", () => {
    updateHistoryBlock();
});

// Monitor changes in the db
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
};

setInterval(checkForUpdates, 1000);

// History panel toggle
let isHistoryOpen = false;

const showHistoryPanel = () => {
    historyPanel.classList.add('open');
    isHistoryOpen = true;
};

const hideHistoryPanel = () => {
    historyPanel.classList.remove('open');
    isHistoryOpen = false;
};

const historyPanelToggle = () => {
    if (isHistoryOpen) {
        hideHistoryPanel();
    } else {
        showHistoryPanel();
    }
};

historyBtn.onclick = () => {
    historyPanelToggle();
};

// Close panel when clicking outside
document.addEventListener('click', (event) => {
    if (isHistoryOpen &&
        !historyPanel.contains(event.target) &&
        !historyBtn.contains(event.target)) {
        hideHistoryPanel();
    }
});

// History item click handlers
document.addEventListener("DOMContentLoaded", () => {
    historyBlock.addEventListener("click", async (event) => {
        if (event.target.classList.contains("history-delete")) {
            const fileId = event.target.dataset.fileId;
            await deleteFileInDB(fileId);
            updateHistoryBlock();
        } else if (event.target.classList.contains("history-go")) {
            const fileId = event.target.dataset.fileId;
            hideHistoryPanel();
            openFile(fileId, true).then(() => {});
        }
    });
});

export {historyPanelToggle};
