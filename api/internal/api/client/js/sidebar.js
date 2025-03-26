import {openFile} from "./app.js";
import {getSortedFilesFromDB, deleteFileInDB} from "./db.js";

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
window.addEventListener("DOMContentLoaded", () => {
    updateHistoryBlock();
});

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
            updateHistoryBlock();
        } else if (event.target.classList.contains("history-go")) {
            const fileId = event.target.dataset.fileId;
            openFile(fileId, true);
        }
    });
});
