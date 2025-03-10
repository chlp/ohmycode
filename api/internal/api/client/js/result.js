const runButton = document.getElementById('run-button');
const cleanResultButton = document.getElementById('clean-result-button');

let runnerBlocksUpdate = undefined;
let resultBlockUpdate = undefined;

window.addEventListener("load", () => {
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
    runnerBlocksUpdate = () => {
        if (file.is_runner_online) {
            runnerContainerBlock.style.display = 'none';
        }
        runnerEditButton.style.display = file.is_runner_online ? 'none' : 'block';
    };

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
    resultBlockUpdate = () => {
        let isRunBtnShouldBeDisabled = false;
        if (file.is_waiting_for_result) {
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

        if (app.actions === 'run' && (file.is_waiting_for_result || file.result.length > 0)) {
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

    let runTask = () => {
        if (!file.is_runner_online) {
            resultCodeMirror.setValue('No runner is available to run your code :(');
            return;
        }
        actions.setContent(contentCodeMirror.getValue());
        file.result = 'In progress..';
        resultCodeMirror.setValue('In progress..');
        runButton.setAttribute('disabled', 'true');
        actions.runTask();
    };
    runButton.onclick = () => {
        runTask();
    };
});