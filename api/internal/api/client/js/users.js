let userOwnNameBlock = document.getElementById('own-name');
let userNameSavingTimeout = null;
let userNameEditing = false;
let userOwnNameEditingFunc = (event) => {
    let key = event.key;
    if (key === 'Backspace' || key === 'Delete' || key === 'ArrowLeft' || key === 'ArrowRight') {
        return true;
    }
    if (key === 'Enter' || key === 'Escape') {
        clearTimeout(userNameSavingTimeout);
        userNameSavingTimeout = null;
        if (userNameEditing) {
            setTimeout(() => {
                updateUsers();
            }, 500);
        }
        userNameEditing = false;
        actions.setUserName();
        event.preventDefault();
        userOwnNameBlock.setAttribute('contenteditable', 'false');
        setTimeout(() => {
            userOwnNameBlock.setAttribute('contenteditable', 'true');
        }, 500);
        contentCodeMirror.focus();
        return false;
    }
    let allowedChars = /^[0-9a-zA-Z_!?:=+\-,.\s'\u0400-\u04ff]*$/;
    if (!allowedChars.test(key)) {
        event.preventDefault();
        return false;
    }
    if (event.target.textContent.length >= 64) {
        event.preventDefault();
        return false;
    }
    userNameEditing = true;
    clearTimeout(userNameSavingTimeout);
    userNameSavingTimeout = setTimeout(() => {
        actions.setUserName();
        if (userNameEditing) {
            setTimeout(() => {
                updateUsers();
            }, 500);
        }
        userNameEditing = false;
    }, 5000);
};

let usersContainerBlock = document.getElementById('users-container');
let usersContainerBlockStateHash = '';
let updateUsers = () => {
    if (userNameEditing || usersContainerBlockStateHash === ohMySimpleHash(file.writer_id + JSON.stringify(file.users))) {
        return;
    }
    usersContainerBlockStateHash = ohMySimpleHash(file.writer_id + JSON.stringify(file.users));
    let users = [];
    let isOwnUserFound = false;
    Object.keys(file.users).forEach((key) => {
        let user = file.users[key];
        user.own = false;
        if (user.id === app.userId) {
            user.own = true;
            app.userName = user.name;
            isOwnUserFound = true;
        }
        users.push(user);
    });
    if (!isOwnUserFound) {
        users.push({
            id: app.userId,
            name: app.userName,
            own: true,
        });
    }
    let html = '';
    if (users.length > 1) {
        users.forEach((user) => {
            if (user.own) {
                html += '<a href="#" id="own-name" contenteditable="true" spellcheck="false" title="Change name">' + user.name + '</a>';
            } else {
                html += '<span>' + user.name + '</span>';
            }
        })
    }
    usersContainerBlock.innerHTML = html;
    if (html !== '') {
        userOwnNameBlock = document.getElementById('own-name');
        userOwnNameBlock.onkeydown = userOwnNameEditingFunc;
    }
};
