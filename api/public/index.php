<!DOCTYPE html>
<html lang="en">
<head>
    <title>OhMyCode</title>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link rel="icon" type="image/png" href="favicon.png">
    <link rel="stylesheet" href="style.css">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.63.0/codemirror.min.css">
    <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.63.0/codemirror.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.63.0/mode/javascript/javascript.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.63.0/mode/go/go.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.63.0/mode/sql/sql.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.63.0/mode/php/php.js"></script>
</head>
<body>

<?php

require __DIR__ . '/bootstrap.php';

$code = '';
$lang = '';
$executor = '';
$executorCheckedAt = null;
$updatedAt = null;

$id = $_GET['session'] ?? null;
if ($id !== null) {
    $id = (string)$id;
    $dbConn = Db::dbConn();
    $stmt = $dbConn->prepare("SELECT `name`, `code`, `lang`, `executor`, `executor_checked_at`, `updated_at`, `writer` FROM `sessions` WHERE `id` = ?");
    if (!$stmt) {
        die('wrong stmt');
    }
    $stmt->bind_param('s', $id);
    $stmt->execute();
    $stmt->bind_result($sessionName, $code, $lang, $executor, $executorCheckedAt, $updatedAt, $writer);
    $stmt->fetch();
    $stmt->close();
}
?>

<div class="blocks-container">
    <button>save</button>
    <input type="text" id="session" style="width: 15em;" maxlength="32" minlength="1"
           pattern="[0-9a-zA-Z\u0400-\u04ff\s\-]{1,32}">
    <label for="session""><- session name</label>
</div>

<div class="blocks-container">
    <button>save</button>
    <input type="text" id="name" style="width: 15em;" maxlength="32" minlength="1"
           pattern="[0-9a-zA-Z\u0400-\u04ff\s\-]{1,32}">
    <label for="name""><- your name (show if not written or clicked change)</label>
</div>

<div class="blocks-container">
    <button>save</button>
    <input type="text" id="executor" style="width: 15em;" maxlength="32" minlength="1"
           pattern="[0-9a-zA-Z]{32}">
    <label for="executor"><- executor (input and hide / show input)</label>
</div>

<div class="blocks-container">
    Session <a href="#">Quinyx 14.11.23</a>
    (<span id="session-status" class="online">online</span>),
    spectators: <span style="">Alex, <a href="#">Serg</a></span>,
    writer: <span style="">Boris</span>
</div>

<div class="editor textarea">
    <textarea id="editor"><?= $code ?></textarea>
</div>
<div class="results textarea">
    <textarea id="results">Waiting for execution...</textarea>
</div>

<div class="blocks-container">
    <button style="float: left; clear: left; margin-right: 1em;">Become a writer</button>
    <select style="width: 120px; float: left; margin-right: 1em;">
        <option>PHP 8.2</option>
        <option>MySQL 8</option>
        <option>GoLang</option>
    </select>
    <button style="float: left;; margin-right: 1em;">Execute code</button>
</div>

<script>
    let user = localStorage['user'];
    if (user === undefined) {
        user = '<?= Utils::genUuid() ?>';
        localStorage['user'] = user;
    }

    window.editor = CodeMirror.fromTextArea(document.getElementById("editor"), {
        lineNumbers: true,
        mode: "sql", // javascript, go, php, sql
        matchBrackets: true,
        indentWithTabs: false,
    });
    // window.editor.setOption('readOnly', true)

    function importCode() {
        var code = window.editor.getValue();
        console.log("Imported Code:", code);
    }

    function updateCode() {
        let scrollInfo = window.editor.getScrollInfo();
        window.editor.setValue("create table sessions\n(\n    id                  varchar(32) not null,\n    code                blob        not null,\n    lang                varchar(32) not null,\n    executor            varchar(32),\n    executor_checked_at datetime,\n    updated_at          datetime(3) default NOW(3) on update NOW(3),\n    constraint sessions_pk\n        primary key (id)\n);\n\ncreate index sessions_executor_idx\n    on sessions (executor);\n\ncreate index sessions_updated_at_idx\n    on sessions (updated_at); create table sessions\n(\n    id                  varchar(32) not null,\n    code                blob        not null,\n    lang                varchar(32) not null,\n    executor            varchar(32),\n    executor_checked_at datetime,\n    updated_at          datetime(3) default NOW(3) on update NOW(3),\n    constraint sessions_pk\n        primary key (id)\n);\n\ncreate index sessions_executor_idx\n    on sessions (executor);\n\ncreate index sessions_updated_at_idx\n    on sessions (updated_at); asd asdasd asdasd");
        window.editor.scrollTo(scrollInfo.left, scrollInfo.top);
    }

    window.results = CodeMirror.fromTextArea(document.getElementById("results"), {
        lineNumbers: true,
        indentWithTabs: false,
        readOnly: true,
    });
</script>

</body>
</html>