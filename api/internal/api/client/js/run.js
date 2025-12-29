import {ohMySimpleHash} from "./utils.js";
import {file} from "./app.js";
import {actions} from "./connect.js";
import {onFileChange} from "./file.js";
import {getCurrentLang, onLangChange} from "./lang.js";
import {contentCodeMirror} from "./editor.js";

const runButton = document.getElementById('run-button');
const cleanResultButton = document.getElementById('clean-result-button');

onLangChange((lang) => {
    switch (lang.action) {
        case 'run':
            runButton.style.display = '';
            cleanResultButton.style.display = '';
            break;
        case 'view':
        case 'edit':
        case 'none':
        default:
            runButton.style.display = 'none';
            cleanResultButton.style.display = 'none';
            break;
    }
});

window.addEventListener("DOMContentLoaded", () => {
    let resultCodeMirror = CodeMirror.fromTextArea(document.getElementById('result'), {
        lineNumbers: true,
        readOnly: true,
        theme: 'tomorrow-night-bright',
    });

    const runnerInput = document.getElementById('runner-input');
    const runnerContainerBlock = document.getElementById('runner-container');
    const fileResultBlock = document.getElementById('file-result');
    const resultContainerBlock = document.getElementById('result-container');

    const runnerSaveButton = document.getElementById('runner-save-button');
    runnerSaveButton.onclick = () => {
        actions.setRunner(runnerInput.value);
    };

    cleanResultButton.onclick = () => {
        file.result = '';
        resultCodeMirror.setValue('');
        actions.cleanResult();
    };

    const runnerEditButton = document.getElementById('runner-edit-button');
    runnerEditButton.onclick = () => {
        runnerEditButtonOnclick();
    };
    const runnerBlocksUpdate = () => {
        if (file.is_runner_online) {
            runnerContainerBlock.style.display = 'none';
        }
        runnerEditButton.style.display = file.is_runner_online ? 'none' : 'block';
    };
    onFileChange(runnerBlocksUpdate);

    let runnerEditButtonOnclick = () => {
        if (runnerContainerBlock.style.display === 'block') {
            runnerContainerBlock.style.display = 'none';
        } else {
            runnerContainerBlock.style.display = 'block';
            runnerInput.focus();
        }
    };
    runnerInput.onkeydown = (event) => {
        if (event.key === 'Enter') {
            event.preventDefault();
            actions.setRunner(runnerInput.value);
        } else if (event.key === 'Escape') {
            runnerEditButtonOnclick();
        }
    };

    let isResultFilledWithInProgress = false;
    // After clicking Run we can force-show "In progress..." for a short window,
    // ignoring any incoming file updates (prevents flicker).
    let forceInProgressUntilMs = 0;
    let forceInProgress = false;
    let runUiToken = 0;
    const resultBlockUpdate = () => {
        const nowMs = Date.now();
        const isForceInProgress = forceInProgress && nowMs < forceInProgressUntilMs;
        const isWaitingForResultUi = file.is_waiting_for_result || isForceInProgress;

        let isRunBtnShouldBeDisabled = false;
        if (isForceInProgress) {
            // Hard lock the UI text to "In progress..." for the whole window,
            // even if server updates try to overwrite it.
            isRunBtnShouldBeDisabled = true;
            isResultFilledWithInProgress = true;
            if (resultCodeMirror.getValue() !== 'In progress...') {
                resultCodeMirror.setValue('In progress...');
            }
        } else if (isWaitingForResultUi) {
            isRunBtnShouldBeDisabled = true;
            if (isResultFilledWithInProgress) {
                resultCodeMirror.setValue(resultCodeMirror.getValue() + '.');
            } else {
                isResultFilledWithInProgress = true;
                resultCodeMirror.setValue('In progress...');
            }
        } else if (file.result.length > 0) {
            if (
                isResultFilledWithInProgress ||
                ohMySimpleHash(file.result) !== ohMySimpleHash(resultCodeMirror.getValue())
            ) {
                isResultFilledWithInProgress = false;
                resultCodeMirror.setValue(file.result);
            }
        } else if (file.is_runner_online) {
            isResultFilledWithInProgress = false;
            resultCodeMirror.setValue('runner will write result here...');
        } else {
            isRunBtnShouldBeDisabled = true;
            isResultFilledWithInProgress = false;
            resultCodeMirror.setValue('...');
        }

        if (isRunBtnShouldBeDisabled) {
            runButton.setAttribute('disabled', 'true');
        } else {
            runButton.removeAttribute('disabled');
        }

        if (getCurrentLang().action === 'run' && (isWaitingForResultUi || file.result.length > 0)) {
            resultContainerBlock.style.display = 'block';
            fileResultBlock.style.display = 'flex';
            cleanResultButton.removeAttribute('disabled');
        } else {
            resultContainerBlock.style.display = 'none';
            fileResultBlock.style.display = 'none';
            cleanResultButton.setAttribute('disabled', 'true');
        }

        resultCodeMirror.refresh();
    };
    onFileChange(resultBlockUpdate);

    let runTask = () => {
        if (!file.is_runner_online) {
            resultCodeMirror.setValue('No runner is available to run your code :(');
            return;
        }
        // Force-show "In progress..." for 1s after clicking Run (UX tweak).
        // During this window we ignore any incoming updates to the result pane,
        // so it can't be overwritten by an older snapshot.
        runUiToken++;
        const myToken = runUiToken;
        forceInProgress = true;
        forceInProgressUntilMs = Date.now() + 1000;
        // Immediately show feedback for the click.
        isResultFilledWithInProgress = true;
        resultCodeMirror.setValue('In progress...');
        // Immediately update layout (button disable / pane visibility) and enforce the lock.
        resultBlockUpdate();
        setTimeout(() => {
            if (runUiToken !== myToken) {
                return;
            }
            forceInProgress = false;
            forceInProgressUntilMs = 0;
            resultBlockUpdate();
        }, 1000);

        // Send content + run as a single atomic action to avoid UI flicker
        // (set_content could otherwise push an update with the previous result).
        runButton.setAttribute('disabled', 'true');
        actions.runTaskWithContent(contentCodeMirror.getValue());
    };
    runButton.onclick = () => {
        runTask();
    };
});