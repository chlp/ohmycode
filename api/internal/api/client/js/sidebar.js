import {openFile} from "./app.js";
import {deleteFileInDB, getSortedFilesFromDB} from "./db.js";

// monitor changes in the db
const updateHistoryBlock = () => {
    getSortedFilesFromDB().then(historyFiles => {
        historyBlock.innerHTML = '';
        historyFiles.forEach((historyFile, index) => {
            if (index > 0) {
                historyBlock.appendChild(document.createElement('br'));
            }
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
if (isSidebarVisible) {
    showSidebar();
} else {
    hideSidebar();
}

const sidebarToggleVisibilitySpan = document.getElementById('sidebar-toggle-visibility');
const sidebarVisibilityToggle = () => {
    if (isSidebarVisible) {
        hideSidebar();
    } else {
        showSidebar();
    }
};
sidebarToggleVisibilitySpan.onclick = () => {
    sidebarVisibilityToggle();
};

document.addEventListener("DOMContentLoaded", () => {
    document.getElementById("history").addEventListener("click", async (event) => {
        if (event.target.classList.contains("history-delete")) {
            const fileId = event.target.dataset.fileId;
            await deleteFileInDB(fileId);
            updateHistoryBlock();
        } else if (event.target.classList.contains("history-go")) {
            const fileId = event.target.dataset.fileId;
            openFile(fileId, true).then(() => {
            });
        }
    });
});

export {sidebarVisibilityToggle};