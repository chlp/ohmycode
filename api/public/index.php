<!DOCTYPE html>
<html lang="en">

<head>
    <title>OhMyCode</title>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link rel="icon" type="image/png" href="favicon.png">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.63.0/codemirror.min.css">
    <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.63.0/codemirror.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.63.0/mode/javascript/javascript.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.63.0/mode/go/go.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.63.0/mode/sql/sql.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.63.0/mode/php/php.js"></script>
    <style>
        .CodeMirror {
            border: 1px solid #666;
            max-width: 1000px;
            width: 50vw;
            min-width: 500px;
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


<textarea id="code"><?= $code ?></textarea>
<button onclick="importCode()">Import Code</button>

<script>
    let editor = CodeMirror.fromTextArea(document.getElementById("code"), {
        lineNumbers: true,
        mode: "sql", // javascript, go, php, sql
        matchBrackets: true,
    });

    function importCode() {
        // Get the code from the CodeMirror editor
        var code = editor.getValue();
        console.log("Imported Code:", code);
        // You can now use the 'code' variable as needed
    }
</script>


</body>
</html>
