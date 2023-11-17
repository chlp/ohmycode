<!DOCTYPE html>
<html lang="en">
<head>
    <title>OhMyCode</title>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link rel="icon" type="image/png" href="favicon.png">
    <link rel="stylesheet" href="style.css">

    <link rel="stylesheet" href="codemirror/codemirror.css">
    <script src="codemirror/codemirror.js"></script>
    <script src="codemirror/mode/clike.js"></script>
    <script src="codemirror/mode/css.js"></script>
    <script src="codemirror/mode/go.js"></script>
    <script src="codemirror/mode/htmlmixed.js"></script>
    <script src="codemirror/mode/javascript.js"></script>
    <script src="codemirror/mode/php.js"></script>
    <script src="codemirror/mode/sql.js"></script>
    <script src="codemirror/mode/xml.js"></script>
</head>
<body>

<?php
require __DIR__ . '/app/bootstrap.php';

$id = trim($_SERVER['REQUEST_URI'], '/');
if (str_contains($id, '?')) {
    $id = substr($id, 0, strpos($id, '?'));
}

$needChangeUrl = false;
if (!Utils::isUuid($id)) {
    $id = Utils::genUuid();
    $needChangeUrl = true;
}

$isNewSession = false;
$session = Session::getById($id);
if ($session === null) {
    $isNewSession = true;
    $session = Session::createNew($id);
}
?>

<div class="blocks-container" id="session-name-container" style="display: none;">
    <button onclick="setSessionName()">save</button>
    <input type="text" id="session-name-input" style="width: 15em;" maxlength="64" minlength="1"
           pattern="[0-9a-zA-Z\u0400-\u04ff\s\-'\.\,]{1,64}">
    <label for="session""><- session name</label>
</div>

<div class="blocks-container" id="user-name-container" style="display: none;">
    <button onclick="setUserName()">save</button>
    <input type="text" id="user-name-input" style="width: 15em;" maxlength="64" minlength="1"
           pattern="[0-9a-zA-Z\u0400-\u04ff\s\-'\.\,]{1,64}">
    <label for="name""><- your name</label>
</div>

<div class="blocks-container">
    Session <a href="#" id="session-name"><?= $session->name ?? '' ?></a><span id="session-status" class="online"></span><span id="users-container"></span>
</div>

<div class="code textarea">
    <textarea id="code"><?= $session->code ?></textarea>
</div>
<div class="result textarea">
    <textarea id="result"><?= $session->result ?? '' ?></textarea>
</div>

<div class="blocks-container">
    <button id="become-writer-button" style="display: none;">Become a writer</button>
    <select id="lang-select" style="width: 150px;">
        <?php
        foreach (Session::LANGS as $key => $data) {
            echo "<option value=\"$key\"";
            if ($session->lang === $key) {
                echo ' selected';
            }
            echo ">{$data['name']}</option>\n";
        }
        ?>
    </select>
    <button id="execute-button" style="display: none">Execute code</button>
    <button onClick="window.open('/', '_blank');" style="float: right;">New session</button>
</div>

<div class="blocks-container" id="executor-container"
     style="float: left; margin-top: 1em; display: <?= $session->isExecutorOnline() ? 'block' : 'none' ?>">
    <button>save</button>
    <input type="text" id="executor" style="width: 20em;" maxlength="32" minlength="32"
           pattern="[0-9a-zA-Z]{32}">
    <label for="executor"><- executor id</label>
</div>

<script>
    <?php
    if ($needChangeUrl) {
        echo "history.pushState({}, null, '/$id');\n";
    }
    ?>
    let initialName = '<?= Utils::randomName() ?>';
    let initialUserId = '<?= Utils::genUuid() ?>';
    let isNewSession = <?= $isNewSession ? 'true' : 'false' ?>;
    let session = <?= $session->getJson() ?>;
    let langKeyToHighlighter = {<?php
        foreach (Session::LANGS as $key => $data) {
            echo "\"$key\": \"{$data['highlighter']}\",";
        }
        ?>};
</script>
<script src="session.js"></script>

</body>
</html>