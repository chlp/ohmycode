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

let codeBlock = CodeMirror.fromTextArea(document.getElementById("code"), {
    lineNumbers: true,
    mode: 'php', // javascript, go, php, sql
    matchBrackets: true,
    indentWithTabs: false,
});
let resultsBlock = CodeMirror.fromTextArea(document.getElementById("results"), {
    lineNumbers: true,
    indentWithTabs: false,
    readOnly: true,
});
// codeBlock.setOption('readOnly', true)
let codeHash = codeBlock.getValue().hash();

let usersContainerBlock = document.getElementById('users-container');
let sessionNameBlock = document.getElementById('session-name');
let sessionStatusBlock = document.getElementById('session-status');

let sessionIsOnline = true;
let ping = undefined;
let userId = localStorage['user'];
if (userId === undefined) {
    userId = initialUserId;
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
        userName = initialName;
        localStorage['tmpUserName'] = userName;
    }
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
    if (isNewSession) {
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
    usersContainerBlock.innerHTML = html;
};
fillUsersContainer();

let postRequest = (url, data, callback, final) => {
    fetch(url, {
        method: 'POST',
        body: JSON.stringify(data),
        headers: {
            "Content-type": "application/json; charset=UTF-8"
        }
    }).then((response) => response.text()).then((text) => callback(text)).finally(() => final());
};

let lastUpdateTimestamp = +new Date;
let pageUpdater = () => {
    let start = +new Date;
    postRequest('/action/session.php', {
        session: session.id,
        user: userId,
        userName: userName,
        action: 'getUpdate',
    }, (response) => {
        ping = +new Date - start;
        console.log((new Date).toLocaleString() + ' | ping: ' + ping);
        lastUpdateTimestamp = +new Date;
        if (response.length === 0) {
            return;
        }
        let data = JSON.parse(response);
        if (data.error !== undefined) {
            console.log(data)
            return
        }
        isNewSession = false;
        session = data;
        if (codeHash !== session.code.hash()) {
            // if writer, not update code
            codeHash = session.code.hash();
            let scrollInfo = codeBlock.getScrollInfo();
            codeBlock.code.setValue(session.code);
            codeBlock.scrollTo(scrollInfo.left, scrollInfo.top);
        }
        fillUsersContainer();
        sessionNameBlock.innerHTML = session.name;
        // set lang: text area, select
        // update session: lang, executorCheckedAt, result, request
    }, () => {
        setTimeout(() => {
            pageUpdater();
        }, 1000);
    });
    if (+new Date - lastUpdateTimestamp > 10000) { // more than 10 seconds
        if (sessionIsOnline) {
            sessionIsOnline = false;
            sessionStatusBlock.classList.remove('online');
            sessionStatusBlock.classList.add('offline');
            sessionStatusBlock.innerHTML = 'offline';
        }
    } else {
        if (!sessionIsOnline) {
            sessionIsOnline = true;
            sessionStatusBlock.classList.remove('offline');
            sessionStatusBlock.classList.add('online');
            sessionStatusBlock.innerHTML = 'online';
        }
    }
};
pageUpdater();