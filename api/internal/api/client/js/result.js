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

cleanResultButton.onclick = () => {
    file.result = '';
    resultCodeMirror.setValue('');
    actions.cleanResult(() => {
        resultBlockUpdate();
    });
};
