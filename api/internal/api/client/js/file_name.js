let fileNameBlock = document.getElementById('file-name');
let fileNameSavingTimeout = null;
let fileNameEditing = false;
fileNameBlock.onkeydown = (event) => {
    let key = event.key;
    if (key === 'Backspace' || key === 'Delete' || key === 'ArrowLeft' || key === 'ArrowRight') {
        return true;
    }
    if (key === 'Enter' || key === 'Escape') {
        clearTimeout(fileNameSavingTimeout);
        fileNameSavingTimeout = null;
        fileNameEditing = false;
        file.name = fileNameBlock.textContent;
        actions.setFileName(fileNameBlock.textContent);
        event.preventDefault();
        fileNameBlock.setAttribute('contenteditable', 'false');
        setTimeout(() => {
            fileNameBlock.setAttribute('contenteditable', 'true');
        }, 500);
        contentCodeMirror.focus();
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
    fileNameEditing = true;
    clearTimeout(fileNameSavingTimeout);
    fileNameSavingTimeout = setTimeout(() => {
        file.name = fileNameBlock.textContent;
        actions.setFileName(fileNameBlock.textContent);
        fileNameEditing = false;
    }, 5000);
};