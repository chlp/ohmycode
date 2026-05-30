const statusBarBlock = document.getElementById('status-bar');

let lockMessage = null;
let tempMessage = null;
let tempClearTimer = null;
let idleMessage = null;

const render = () => {
    const text = lockMessage ?? tempMessage ?? idleMessage ?? null;
    if (text === null) {
        statusBarBlock.style.display = 'none';
        statusBarBlock.textContent = '';
        return;
    }
    statusBarBlock.textContent = text;
    statusBarBlock.style.removeProperty('display');
};

// High priority: Offline / Blocked by another user. Permanent until cleared.
const setLockStatus = (text) => {
    lockMessage = text || null;
    render();
};

// Medium priority: transient notification (run time, save confirmation).
// Clears after clearAfterMs if > 0.
const setStatus = (text, clearAfterMs) => {
    if (tempClearTimer) {
        clearTimeout(tempClearTimer);
        tempClearTimer = null;
    }
    tempMessage = text || null;
    if (clearAfterMs > 0 && tempMessage !== null) {
        tempClearTimer = setTimeout(() => {
            tempMessage = null;
            tempClearTimer = null;
            render();
        }, clearAfterMs);
    }
    render();
};

// Low priority: always-visible background info (content size).
const setIdleStatus = (text) => {
    idleMessage = text || null;
    render();
};

export {setLockStatus, setStatus, setIdleStatus};
