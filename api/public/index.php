<!DOCTYPE html>
<html lang="en">
<head>
    <title>OhMyCode</title>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">

    <link rel="icon" type="image/svg+xml" href="favicon.svg">
    <link rel="icon" type="image/png" href="favicon.png">
    <link rel="icon" type="image/png" href="favicon-96x96.png" sizes="96x96" />
    <link rel="shortcut icon" href="favicon.ico" />
    <link rel="apple-touch-icon" sizes="180x180" href="favicon-apple-touch.png" />

    <link rel="stylesheet" href="style.css?v=1">

    <link rel="stylesheet" href="codemirror/codemirror.css">
    <link rel="stylesheet" href="codemirror/themes/base16-light.css">
    <link rel="stylesheet" href="codemirror/themes/base16-dark.css">
    <link rel="stylesheet" href="codemirror/themes/tomorrow-night-bright.css">
    <script src="codemirror/codemirror.js"></script>
    <script src="codemirror/mode/clike.js"></script>
    <script src="codemirror/mode/css.js"></script>
    <script src="codemirror/mode/go.js"></script>
    <script src="codemirror/mode/htmlmixed.js"></script>
    <script src="codemirror/mode/javascript.js"></script>
    <script src="codemirror/mode/markdown.js"></script>
    <script src="codemirror/mode/php.js"></script>
    <script src="codemirror/mode/sql.js"></script>
    <script src="codemirror/mode/xml.js"></script>

    <script src="js/utils.js?v=1"></script>
    <script>
        let sessionId = window.location.pathname.slice(1);
        if (!isUuid(sessionId)) {
            sessionId = genUuid();
            history.pushState({}, null, '/' + sessionId);
        }
    </script>
</head>
<body>

<div class="blocks-container" id="session-name-container" style="float: left; clear: left;">
    <a href="#" id="session-name" contenteditable="true" spellcheck="false"
       title="Rename file"></a><span id="session-status" class="online"></span>
</div>

<div class="blocks-container" style="float: right; clear: right;">
    <span id="users-container"></span>
    <a href="https://github.com/chlp/ohmycode" target="_blank">
        <img src="github-mark.svg" style="height: 28px; vertical-align: middle; margin-left: 1em;">
    </a>
</div>

<div class="code textarea" id="code-container" style="clear: both;">
    <textarea id="code"></textarea>
</div>

<div class="blocks-container" id="controls-container" style="float: left; clear: left; display: none;">
    <select id="lang-select" style="width: 150px; height: 30px;"></select>
    <button id="run-button" title="Cmd/Ctrl + Enter" disabled>Run code</button>
    <button id="clean-result-button" disabled>Clean result</button>
    <button onclick="runnerEditButtonOnclick()" id="runner-edit-button" style="display: none;">Runner</button>
    <button onclick="copyToClipboard(window.location.href)">Copy URL</button>
    <a href="/" class="button" target="_blank">New file</a>
</div>

<div class="blocks-container" style="float: right; clear: right; padding: 2px 0;">
    <span id="current-writer-info" style="padding: 0.4rem 0.8rem; display: none;">
        Code is writing now by <span
                id="current-writer-name"></span>
    </span>
</div>

<div class="blocks-container" id="runner-container" style="float: left; margin-top: 1em; display: none;">
    <button id="runner-save-button">save</button>
    <input type="text" id="runner-input" style="width: 20em;" maxlength="32" minlength="32"
           pattern="[0-9a-zA-Z]{32}" value="">
    <label for="runner-input"><- runner id</label>
</div>

<div class="result textarea" id="result-container">
    <textarea id="result"></textarea>
</div>
<script src="js/actions.js?v=1"></script>
<script src="js/session.js?v=1"></script>
<script src="js/session_name.js?v=1"></script>
<script src="js/users.js?v=1"></script>

</body>
</html>
