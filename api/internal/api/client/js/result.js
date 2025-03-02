const runnerInput = document.getElementById('runner-input');
const runButton = document.getElementById('run-button');
const runnerContainerBlock = document.getElementById('runner-container');
const resultContainerBlock = document.getElementById('result-container');

const runnerSaveButton = document.getElementById('runner-save-button');
runnerSaveButton.onclick = () => {
    actions.setRunner();
};


const cleanResultButton = document.getElementById('clean-result-button');
cleanResultButton.onclick = () => {
    file.result = '';
    resultCodeMirror.setValue('');
    actions.cleanResult(() => {
        resultBlockUpdate();
    });
};

const runnerEditButton = document.getElementById('runner-edit-button');
let runnerBlocksUpdate = () => {
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
        actions.setRunner();
    } else if (event.key === 'Escape') {
        runnerEditButtonOnclick();
    }
};

let resultCodeMirror = CodeMirror.fromTextArea(document.getElementById('result'), {
    lineNumbers: true,
    readOnly: true,
    theme: 'tomorrow-night-bright',
});
let isResultFilledWithInProgress = false;
let resultBlockUpdate = () => {
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

    if (file.is_waiting_for_result || file.result.length > 0) {
        resultContainerBlock.style.display = 'block';
        contentContainerBlock.style.height = 'calc(68vh - 90px)';
        cleanResultButton.removeAttribute('disabled');
    } else {
        resultContainerBlock.style.display = 'none';
        contentContainerBlock.style.height = 'calc(98vh - 90px)';
        cleanResultButton.setAttribute('disabled', 'true');
    }
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
contentContainerBlock.onkeydown = (event) => {
    if ((event.metaKey || event.ctrlKey) && event.key === 'Enter') {
        runTask();
    }
};
