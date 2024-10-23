let sessionNameBlock = document.getElementById('session-name');
let sessionNameSavingTimeout = null;
let sessionNameEditing = false;
sessionNameBlock.onkeydown = (event) => {
    let key = event.key;
    if (key === 'Backspace' || key === 'Delete' || key === 'ArrowLeft' || key === 'ArrowRight') {
        return true;
    }
    if (key === 'Enter' || key === 'Escape') {
        clearTimeout(sessionNameSavingTimeout);
        sessionNameSavingTimeout = null;
        sessionNameEditing = false;
        actions.setSessionName();
        event.preventDefault();
        sessionNameBlock.setAttribute('contenteditable', 'false');
        setTimeout(() => {
            sessionNameBlock.setAttribute('contenteditable', 'true');
        }, 500);
        contentBlock.focus();
        return false;
    }
    let allowedChars = /^[0-9a-zA-Z_!?:=+\-,.\s'\u0400-\u04ff]*$/;
    if (!allowedChars.test(key)) {
        event.preventDefault();
        return false;
    }
    if (event.target.textContent.length >= 64) {
        event.preventDefault();
        return false;
    }
    sessionNameEditing = true;
    clearTimeout(sessionNameSavingTimeout);
    sessionNameSavingTimeout = setTimeout(() => {
        actions.setSessionName();
        sessionNameEditing = false;
    }, 5000);
};