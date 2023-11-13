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
    <style>
        .CodeMirror {
            border: 1px solid #666;
            width: 50vw;
            min-width: 500px;
            height: 80vh;
            min-height: 300px;
            cursor: text;
            float: left;
        }
    </style>
</head>
<body>

<?php
$code = '';
$lang = '';
$executor = '';
$executorCheckedAt = null;
$updatedAt = null;

$id = $_GET['session'] ?? null;
if ($id !== null) {
    $id = (string)$id;
    require 'db.php';
    $dbConn = dbConn();
    $stmt = $dbConn->prepare("SELECT `code`, `lang`, `executor`, `executor_checked_at`, `updated_at` FROM `sessions` WHERE `id` = ?");
    if (!$stmt) {
        die('wrong stmt');
    }
    $stmt->bind_param('s', $id);
    $stmt->execute();
    $stmt->bind_result($code, $lang, $executor, $executorCheckedAt, $updatedAt);
    $stmt->fetch();
    $stmt->close();
}
?>

<div class="header">
    <button style="float: left; clear: left; margin-right: 1em;">save</button>
    <input type="text" id="name" style="width: 15em; float: left;">
    <label for="name" style="float: left;"><- your name (show if not written or clicked change)</label>

    <button style="float: left; clear: left; margin-right: 1em;">save</button>
    <input type="text" id="executor" style="width: 15em; float: left;">
    <label for="executor" style="float: left;"><- executor (input and hide / show input)</label>

    <div style="clear: both;">
        Session <span style="color: cornflowerblue">Quinyx 14.11.23</span>,
        executor: <span style="color: forestgreen;">online</span>,
        spectators: <span style="">Alex, <u>Serg</u></span>,
        writer: <span style="">Boris</span>
    </div>
</div>

<textarea id="code"><?= $code ?></textarea>
<div style="float:left; border: 1px solid #666; width: 40vw; height: 80vh; margin-left: 5px;">results</div>

<div style="float: left; clear: both; padding: 1em;">
    <input type="button" value="Become a writer" style="float: left;">
    <select style="width: 120px; float: left;">
        <option>PHP 8.2</option>
        <option>MySQL 8</option>
        <option>GoLang</option>
    </select>
    <input type="button" value="Execute code"><br><br>
</div>

<script>
    window.editor = CodeMirror.fromTextArea(document.getElementById("code"), {
        lineNumbers: true,
        mode: "sql", // javascript, go, php, sql
        matchBrackets: true,
        indentWithTabs: false,
    });

    function importCode() {
        var code = window.editor.getValue();
        console.log("Imported Code:", code);
    }

    function updateCode() {
        let scrollInfo = window.editor.getScrollInfo();
        window.editor.setValue("create table sessions\n(\n    id                  varchar(32) not null,\n    code                blob        not null,\n    lang                varchar(32) not null,\n    executor            varchar(32),\n    executor_checked_at datetime,\n    updated_at          datetime(3) default NOW(3) on update NOW(3),\n    constraint sessions_pk\n        primary key (id)\n);\n\ncreate index sessions_executor_idx\n    on sessions (executor);\n\ncreate index sessions_updated_at_idx\n    on sessions (updated_at); create table sessions\n(\n    id                  varchar(32) not null,\n    code                blob        not null,\n    lang                varchar(32) not null,\n    executor            varchar(32),\n    executor_checked_at datetime,\n    updated_at          datetime(3) default NOW(3) on update NOW(3),\n    constraint sessions_pk\n        primary key (id)\n);\n\ncreate index sessions_executor_idx\n    on sessions (executor);\n\ncreate index sessions_updated_at_idx\n    on sessions (updated_at); asd asdasd asdasd");
        window.editor.scrollTo(scrollInfo.left, scrollInfo.top);
    }
</script>


</body>
</html>
