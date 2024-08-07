let userOwnNameBlock = null;
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
        userNameEditing = false;
        actions.setUserName();
        event.preventDefault();
        userOwnNameBlock.setAttribute('contenteditable', 'false');
        setTimeout(() => {
            userOwnNameBlock.setAttribute('contenteditable', 'true');
        }, 500);
        codeBlock.focus();
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
        userNameEditing = false;
    }, 5000);
};

let usersContainerBlock = document.getElementById('users-container');
let usersContainerState = '';
let updateUsers = () => {
    if (sessionPreviousState.writer + JSON.stringify(sessionPreviousState.users) === session.writer + JSON.stringify(session.users)) {
        return;
    }
    let spectators = [];
    let writer = undefined;
    if (isNewSession) {
        isWriter = true;
        writer = {
            id: userId,
            name: userName,
            own: true,
        };
        spectators = [];
    } else {
        isWriter = userId === session.writer;
        session.users.forEach((user) => {
            user.own = false;
            if (user.id === userId) {
                user.own = true;
                userName = user.name;
            }
            if (user.id === session.writer) {
                writer = user;
            } else {
                spectators.push(user);
            }
        });
    }
    let html = '';
    if (writer !== undefined) {
        if (writer.own) {
            html += '<a href="#" id="own-name" contenteditable="true" spellcheck="false" class="writer">' + writer.name + '</a>';
        } else {
            html += '<span class="writer">' + writer.name + '</span>';
        }
    }
    if (spectators.length > 0) {
        spectators.forEach((user, i) => {
            if (user.own) {
                html += '<a href="#" id="own-name" contenteditable="true" spellcheck="false">' + user.name + '</a>';
            } else {
                html += '<span>' + user.name + '</span>';
            }
        })
    }
    if (usersContainerState !== html.hash() && !userNameEditing) {
        usersContainerState = html.hash();
        usersContainerBlock.innerHTML = html;
        userOwnNameBlock = document.getElementById('own-name');
        userOwnNameBlock.onkeydown = userOwnNameEditingFunc;
    }
};
document.addEventListener('DOMContentLoaded', () => {
    updateUsers();
});