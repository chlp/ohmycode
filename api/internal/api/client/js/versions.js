import {actions, onVersions, onOpenFile} from "./connect.js";

const versionsPanel = document.getElementById('versions-panel');
const versionsBlock = document.getElementById('versions');
const versionsBtn = document.getElementById('header-versions-btn');

let isVersionsOpen = false;
let versionsList = [];
let pendingNewTab = null;

// Register versions handler
onVersions((versions) => {
    setVersions(versions);
});

// Register open_file handler - redirect pending tab to new file
onOpenFile((fileId) => {
    console.log('onOpenFile called with fileId:', fileId, 'pendingNewTab:', pendingNewTab);
    if (pendingNewTab && !pendingNewTab.closed) {
        pendingNewTab.location.href = '/' + fileId;
        pendingNewTab = null;
    }
});

const formatDate = (dateStr) => {
    const date = new Date(dateStr);
    return new Intl.DateTimeFormat('sv-SE', {
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
    }).format(date).replace(',', '');
};

const updateVersionsBlock = () => {
    versionsBlock.innerHTML = '';

    if (versionsList.length === 0) {
        let emptyMsg = document.createElement('div');
        emptyMsg.className = 'versions-empty';
        emptyMsg.textContent = 'No versions yet';
        versionsBlock.appendChild(emptyMsg);
        return;
    }

    versionsList.forEach((version) => {
        let span = document.createElement('span');
        span.className = 'version-item';

        let restoreLink = document.createElement('a');
        restoreLink.className = 'version-restore';
        restoreLink.dataset.versionId = version.id;
        restoreLink.textContent = version.name;
        restoreLink.title = `${version.lang}`;

        let dateSpan = document.createElement('span');
        dateSpan.className = 'version-date';
        dateSpan.textContent = formatDate(version.created_at);

        span.appendChild(restoreLink);
        span.appendChild(dateSpan);
        versionsBlock.appendChild(span);
    });
};

const setVersions = (versions) => {
    versionsList = versions || [];
    updateVersionsBlock();
};

const showVersionsPanel = () => {
    versionsPanel.classList.add('open');
    isVersionsOpen = true;
    actions.getVersions();
};

const hideVersionsPanel = () => {
    versionsPanel.classList.remove('open');
    isVersionsOpen = false;
};

const versionsPanelToggle = () => {
    if (isVersionsOpen) {
        hideVersionsPanel();
    } else {
        showVersionsPanel();
    }
};

versionsBtn.onclick = () => {
    versionsPanelToggle();
};

// Close panel when clicking outside
document.addEventListener('click', (event) => {
    if (isVersionsOpen &&
        !versionsPanel.contains(event.target) &&
        !versionsBtn.contains(event.target)) {
        hideVersionsPanel();
    }
});

// Version item click handlers
document.addEventListener("DOMContentLoaded", () => {
    versionsBlock.addEventListener("click", async (event) => {
        if (event.target.classList.contains("version-restore")) {
            const versionId = event.target.dataset.versionId;
            // Open new tab immediately (must be in click handler to avoid popup blocker)
            pendingNewTab = window.open('about:blank', '_blank');
            actions.restoreVersion(versionId);
            hideVersionsPanel();
        }
    });
});

export {versionsPanelToggle};
