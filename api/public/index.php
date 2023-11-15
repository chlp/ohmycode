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
    <input type="text" id="session-name-input" style="width: 15em;" maxlength="32" minlength="1"
           pattern="[0-9a-zA-Z\u0400-\u04ff\s\-\']{1,32}">
    <label for="session""><- session name</label>
</div>

<div class="blocks-container" id="user-name-container" style="display: none;">
    <button>save</button>
    <input type="text" id="user-name-input" style="width: 15em;" maxlength="32" minlength="1"
           pattern="[0-9a-zA-Z\u0400-\u04ff\s\-\']{1,32}">
    <label for="name""><- your name</label>
</div>

<div class="blocks-container">
    Session <a href="#" id="session-name"><?= $session->name ?? '' ?></a>
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

    let session = <?= $session->getJson() ?>;
    let sessionId = '<?= $session->id ?>';
    let userId = localStorage['user'];
    if (userId === undefined) {
        userId = '<?= Utils::genUuid() ?>';
        localStorage['user'] = userId;
    }
    let userName = undefined;
    session.users.forEach((user) => {
        if (user.id === userId) {
            userName = user.name;
        }
    });
    if (userName === undefined) {
        userName = localStorage['tmpUserName'];
        if (userName === undefined) {
            userName = '<?= Utils::randomName() ?>';
            localStorage['tmpUserName'] = userName;
        }
    }
    let newSession = <?= $newSession ? 'true' : 'false' ?>;
    let codeHash = undefined;

    String.prototype.hash = function () {
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

    let fillUsersContainer = () => {
        let spectators = [];
        let writer = undefined;
        session.users.forEach((user) => {
            user.own = user.id === userId;
            if (user.id === session.writer) {
                writer = user;
            } else {
                spectators.push(user)
            }
        });
        if (newSession) {
            writer = {
                id: userId,
                name: userName,
                own: true,
            }
        }
        let html = '';
        if (writer !== undefined) {
            html += ', writer: ';
            if (writer.own) {
                html += '<a id="own-name" href="#">';
            }
            html += writer.name;
            if (writer.own) {
                html += '</a>';
            }
        }
        if (spectators.length > 0) {
            html += ', spectators: ';
            spectators.forEach((user) => {
                if (user.own) {
                    html += '<a id="own-name" href="#">';
                }
                html += user.name;
                if (user.own) {
                    html += '</a>';
                }
            })
        }
        document.getElementById('users-container').innerHTML = html;
    };
    fillUsersContainer();

    let lastUpdateTimestamp = +new Date / 1000;
    setInterval(() => {
        let start = parseInt(+new Date);
        postRequest('/action/session.php', {
            session: sessionId,
            user: userId,
            userName: userName,
            action: 'getUpdate',
        }, (response) => {
            let ping = parseInt(+new Date) - start;
            console.log('ping ' + ping);
            lastUpdateTimestamp = +new Date / 1000;
            if (response.length === 0) {
                return;
            }
            let data = JSON.parse(response);
            if (data.error !== undefined) {
                console.log(data)
                return
            }
            newSession = false;
            session = data;
            if (codeHash !== session.code.hash()) {
                codeHash = session.code.hash();
                let scrollInfo = window.code.getScrollInfo();
                window.code.setValue(session.code);
                window.code.scrollTo(scrollInfo.left, scrollInfo.top);
            }
            fillUsersContainer();
            document.getElementById('session-name').innerHTML = session.name;
            // set lang: text area, select
            // update session: lang, executorCheckedAt, result, request
            // set online ping
        });
        if (+new Date / 1000 - lastUpdateTimestamp > 10) {
            let sessionStatusEl = document.getElementById('session-status');
            sessionStatusEl.classList.remove('online');
            sessionStatusEl.classList.add('offline');
            sessionStatusEl.innerHTML = 'offline';
        } else {
            let sessionStatusEl = document.getElementById('session-status');
            sessionStatusEl.classList.remove('offline');
            sessionStatusEl.classList.add('online');
            sessionStatusEl.innerHTML = 'online';
        }
    }, 1000);

    let postRequest = (url, data, callback) => {
        fetch(url, {
            method: 'POST',
            body: JSON.stringify(data),
            headers: {
                "Content-type": "application/json; charset=UTF-8"
            }
        }).then((response) => response.text()).then((text) => callback(text));
    };
</script>

</body>
</html>