import {ohMySimpleHash} from "./utils.js";
import {app, file, openFile} from "./app.js";
import {saveFileToDB} from "./db.js";
import {contentCodeMirror, contentMarkdownBlock} from "./editor.js";

const fileChangeHandlers = [];
const onFileChange = (callback) => {
    if (typeof callback === "function") {
        fileChangeHandlers.push(callback);
    }
};

const controlsContainerBlock = document.getElementById('controls-container');
const applyFile = (newFile) => {
    if (file.id !== newFile.id) {
        // Ignore stale WS updates for a different file (can happen during navigation/reconnect).
        // Re-opening the current file here causes UI reset loops and breaks the editor on load.
        console.log('onmessage: ignore other file.id', newFile.id, 'current=', file.id);
        return;
    }

    let previousWriterId = file.writer_id;

    file.name = newFile.name;
    file.lang = newFile.lang;
    file.runner = newFile.runner;
    file.is_runner_online = newFile.is_runner_online;
    file.updated_at = newFile.updated_at;
    file.content_updated_at = newFile.content_updated_at;
    file.users = newFile.users;
    file.is_waiting_for_result = newFile.is_waiting_for_result;
    file.result = newFile.result;
    file.persisted = newFile.persisted;
    file.writer_id = newFile.writer_id;
    if (typeof newFile.content === 'string') {
        file.content = newFile.content;
    }

    if (file.persisted) {
        saveFileToDB(file).then(() => {
        });
    }
    document.title = `OhMyCode – ${file.name}`;

    fileChangeHandlers.forEach(fn => fn(file));

    if (
        !app.isOnline || // first load
        (
            file.writer_id !== app.id && previousWriterId !== app.id && // do not update if current user is writer
            ohMySimpleHash(file.content) !== ohMySimpleHash(contentCodeMirror.getValue()) // do not update if code is the same already
        )
    ) {
        let {left, top} = contentCodeMirror.getScrollInfo();
        let {line, ch} = contentCodeMirror.getCursor();
        contentCodeMirror.setValue(file.content);
        contentMarkdownBlock.innerHTML = marked.parse(file.content);
        contentCodeMirror.scrollTo(left, top);
        contentCodeMirror.setCursor({line: line, ch: ch});
    }

    if (controlsContainerBlock.style.display !== 'block') {
        controlsContainerBlock.style.display = 'block';
    }
};

export {applyFile, onFileChange};