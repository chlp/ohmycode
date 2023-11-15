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

$newSession = false;
$session = Session::getById($id);
if ($session === null) {
    $newSession = true;
    $session = Session::createNew($id);
}
?>

<div class="blocks-container" id="session-name-container" style="display: none;">
    <button>save</button>
    <input type="text" id="session-name" style="width: 15em;" maxlength="32" minlength="1"
           pattern="[0-9a-zA-Z\u0400-\u04ff\s\-\']{1,32}">
    <label for="session""><- session name</label>
</div>

<div class="blocks-container" id="user-name-container" style="display: none;">
    <button>save</button>
    <input type="text" id="user-name" style="width: 15em;" maxlength="32" minlength="1"
           pattern="[0-9a-zA-Z\u0400-\u04ff\s\-\']{1,32}">
    <label for="name""><- your name</label>
</div>

<div class="blocks-container">
    Session <a href="#"><?= $session->name ?? '' ?></a>
    (<span id="session-status" class="online">online</span>)<span id="users-container"></span>
</div>

<div class="code textarea">
    <textarea id="code"><?= $session->code ?? '' ?></textarea>
</div>
<div class="results textarea">
    <textarea id="results">Waiting for execution...</textarea>
</div>

<div class="blocks-container">
    <button>Become a writer</button>
    <select id="lang" style="width: 150px;">
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
    <button onClick="window.open('/', '_blank');">New session</button>
</div>

<div class="blocks-container">
    <button>save</button>
    <input type="text" id="executor" style="width: 15em;" maxlength="32" minlength="1"
           pattern="[0-9a-zA-Z]{32}">
    <label for="executor"><- executor id</label>
</div>

<script>
    <?php
    if ($needChangeUrl) {
        echo "history.pushState({}, null, '/$id');\n";
    }
    ?>

    let userId = localStorage['user'];
    let tmpUserName = '<?= Utils::randomName() ?>';
    if (userId === undefined) {
        userId = '<?= Utils::genUuid() ?>';
        localStorage['user'] = userId;
    }
    let sessionUpdatedAt = '<?= $session->updatedAt->format('Y-m-d H:i:s.u') ?>';
    let newSession = <?= $newSession ? 'true' : 'false' ?>;

    String.prototype.hashCode = function () {
        let hash = 0,
            i, chr;
        if (this.length === 0) return hash;
        for (i = 0; i < this.length; i++) {
            chr = this.charCodeAt(i);
            hash = ((hash << 5) - hash) + chr;
            hash |= 0;
        }
        return hash;
    };

    window.code = CodeMirror.fromTextArea(document.getElementById("code"), {
        lineNumbers: true,
        mode: 'php', // javascript, go, php, sql
        matchBrackets: true,
        indentWithTabs: false,
    });
    window.results = CodeMirror.fromTextArea(document.getElementById("results"), {
        lineNumbers: true,
        indentWithTabs: false,
        readOnly: true,
    });

    // window.code.setOption('readOnly', true)

    function importCode() {
        let code = window.code.getValue();
        console.log("Imported Code:", code);
    }

    let getUsersContainerContent = () => {
        if (newSession) {
            return ', writer: <a id="own-name" href="#">' + tmpUserName + '</a>';
        }
        let html = ', spectators: <a id="own-name" href="#">' + tmpUserName + '</a>';
        return html;
    };
    let fillUserContainer = () => {
        document.getElementById('users-container').innerHTML = getUsersContainerContent();
    }
    fillUserContainer();

    let updateCode = () => {
        // receive session: name, code, lang, users, writer, updatedAt, executorCheckedAt, result, request
        // calc ping
        let scrollInfo = window.code.getScrollInfo();
        window.code.setValue("create table sessions\n(\n    id                  varchar(32) not null,\n    code                blob        not null,\n    lang                varchar(32) not null,\n    executor            varchar(32),\n    executor_checked_at datetime,\n    updated_at          datetime(3) default NOW(3) on update NOW(3),\n    constraint sessions_pk\n        primary key (id)\n);\n\ncreate index sessions_executor_idx\n    on sessions (executor);\n\ncreate index sessions_updated_at_idx\n    on sessions (updated_at); create table sessions\n(\n    id                  varchar(32) not null,\n    code                blob        not null,\n    lang                varchar(32) not null,\n    executor            varchar(32),\n    executor_checked_at datetime,\n    updated_at          datetime(3) default NOW(3) on update NOW(3),\n    constraint sessions_pk\n        primary key (id)\n);\n\ncreate index sessions_executor_idx\n    on sessions (executor);\n\ncreate index sessions_updated_at_idx\n    on sessions (updated_at); asd asdasd asdasd");
        window.code.scrollTo(scrollInfo.left, scrollInfo.top);
    }
    setInterval(() => {
        updateCode();
    }, 1000);
</script>

</body>
</html>