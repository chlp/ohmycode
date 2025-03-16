import {app, file} from "./app.js";
import {actions} from "./connect.js";
import {contentCodeMirror, contentMarkdownBlock} from "./editor.js";
import {fileNameBlock} from "./file_name.js";

contentCodeMirror.on('drop', (cm, event) => {
    event.preventDefault();
});

document.addEventListener('dragover', (event) => {
    event.preventDefault();
});

document.addEventListener('drop', (event) => {
    event.preventDefault();
    const droppedFiles = event.dataTransfer.files;
    if (droppedFiles.length === 0) {
        return;
    }
    const droppedFile = droppedFiles[0];
    if (droppedFile.size > 512 * 1024) {
        console.warn('File too large (>512Kb)', droppedFile);
        return;
    }
    const reader = new FileReader();
    reader.onload = async (e) => {
        if (await isFileBinary(droppedFile)) {
            console.warn("Wrong file (binary)", droppedFile);
            return;
        }
        if (file.writer_id !== '' && file.writer_id !== app.id) {
            return;
        }

        let newFileName = droppedFile.name;
        let newContent = e.target.result;
        const allowedCharsRegex = /[^0-9a-zA-Z_!?:=+\-,.\s'\u0400-\u04ff]/g;
        newFileName = newFileName.replace(allowedCharsRegex, '');
        newFileName = newFileName.substring(0, 64);

        fileNameBlock.innerHTML = newFileName;
        file.name = newFileName;
        actions.setFileName(newFileName);

        contentCodeMirror.setValue(newContent);
        contentMarkdownBlock.innerHTML = marked.parse(file.content);
        actions.setContent(newContent);
    };
    reader.onerror = function () {
        console.error('Error occurred: ' + droppedFile);
    };
    reader.readAsText(droppedFile);
});

let isFileBinary = async (file) => {
    const buffer = await file.arrayBuffer();
    const bytes = new Uint8Array(buffer);
    const maxBytesToCheck = Math.min(bytes.length, 32 * 1024);
    let nonPrintableCount = 0;
    for (let i = 0; i < maxBytesToCheck; i++) {
        const byte = bytes[i];
        if ((byte < 32 || byte > 126) && byte !== 9 && byte !== 10 && byte !== 13) {
            nonPrintableCount++;
        }
    }
    let nonPrintableRateForBinary = 0.4;
    if (bytes.length < 1024) {
        nonPrintableRateForBinary = 0.6;
    }
    return nonPrintableCount / maxBytesToCheck > nonPrintableRateForBinary;
}
