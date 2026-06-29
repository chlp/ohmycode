import {generateKey, importKey} from "./encrypt.js";
import {app, file} from "./app.js";
import {actions} from "./connect.js";
import {onFileChange} from "./file.js";
import {contentCodeMirror} from "./editor.js";

const encryptPanel = document.getElementById('encrypt-panel');
const encryptPanelContent = document.getElementById('encrypt-panel-content');
const encryptBtn = document.getElementById('header-encrypt-btn');
const noKeyOverlay = document.getElementById('no-key-overlay');
const noKeyInput = document.getElementById('no-key-input');
const noKeyError = document.getElementById('no-key-error');

let isEncryptPanelOpen = false;

const showEncryptPanel = () => {
    encryptPanel.classList.add('open');
    isEncryptPanelOpen = true;
    renderEncryptPanel();
};

const hideEncryptPanel = () => {
    encryptPanel.classList.remove('open');
    isEncryptPanelOpen = false;
};

const encryptPanelToggle = () => {
    if (isEncryptPanelOpen) {
        hideEncryptPanel();
    } else {
        showEncryptPanel();
    }
};

const copyText = async (text, btn) => {
    try {
        await navigator.clipboard.writeText(text);
        const orig = btn.textContent;
        btn.textContent = 'Copied!';
        setTimeout(() => { btn.textContent = orig; }, 1500);
    } catch(e) {
        console.error('Copy failed:', e);
    }
};

const buildShareLinks = () => {
    const keyStr = localStorage.getItem('ohmycode_key_' + file.id) || '';
    const roKeyStr = localStorage.getItem('ohmycode_rokey_' + file.id) || '';
    const base = window.location.origin + '/' + file.id;
    const editLink = keyStr ? base + '#key=' + keyStr : base;
    let roLink = null;
    if (file.ro_token) {
        const roHash = roKeyStr
            ? '#key=' + roKeyStr + '&ro=' + file.ro_token
            : '#ro=' + file.ro_token;
        roLink = base + roHash;
    }
    return { editLink, roLink, editKey: keyStr, roKey: roKeyStr };
};

const addLinkRow = (label, url) => {
    const row = document.createElement('div');
    row.className = 'encrypt-link-row';

    const lbl = document.createElement('div');
    lbl.className = 'encrypt-link-label';
    lbl.textContent = label;

    const inputWrap = document.createElement('div');
    inputWrap.className = 'encrypt-link-input-wrap';

    const input = document.createElement('input');
    input.type = 'text';
    input.readOnly = true;
    input.value = url;
    input.className = 'encrypt-link-input';
    input.onclick = () => input.select();

    const btn = document.createElement('button');
    btn.textContent = 'Copy';
    btn.className = 'encrypt-copy-btn';
    btn.onclick = () => copyText(url, btn);

    inputWrap.appendChild(input);
    inputWrap.appendChild(btn);
    row.appendChild(lbl);
    row.appendChild(inputWrap);
    encryptPanelContent.appendChild(row);
};

const renderEncryptPanel = () => {
    encryptPanelContent.innerHTML = '';

    if (app.isROLink) {
        const msg = document.createElement('p');
        msg.className = 'encrypt-panel-msg';
        msg.textContent = file.encrypted
            ? 'Read-only access. You cannot edit or change encryption settings.'
            : 'You have read-only access via a shared link.';
        encryptPanelContent.appendChild(msg);

        if (file.encrypted) {
            const roKeyStr = localStorage.getItem('ohmycode_rokey_' + file.id) || '';
            const base = window.location.origin + '/' + file.id;
            if (roKeyStr) addLinkRow('Read-only key:', roKeyStr);
            if (app.roToken) {
                const roLink = roKeyStr
                    ? base + '#key=' + roKeyStr + '&ro=' + app.roToken
                    : base + '#ro=' + app.roToken;
                addLinkRow('Read-only link:', roLink);
            }
        }
        return;
    }

    if (!file.encrypted) {
        const msg = document.createElement('p');
        msg.className = 'encrypt-panel-msg';
        msg.textContent = 'Enable encryption to keep content private. The key is stored in your browser and never sent to the server.';
        encryptPanelContent.appendChild(msg);

        const btn = document.createElement('button');
        btn.textContent = 'Enable Encryption';
        btn.className = 'encrypt-action-btn';
        btn.onclick = async (e) => {
            e.stopPropagation();
            btn.disabled = true;
            btn.textContent = 'Generating key…';
            try {
                const {key: editKey, exported: editExported} = await generateKey();
                const {key: roKey, exported: roExported} = await generateKey();
                app.encKey = editKey;
                app.roEncKey = roKey;
                localStorage.setItem('ohmycode_key_' + file.id, editExported);
                localStorage.setItem('ohmycode_rokey_' + file.id, roExported);
                file.encrypted = true;
                actions.setEncrypted(true);
                actions.setContent(contentCodeMirror.getValue());
                renderEncryptPanel();
            } catch(e) {
                console.error('Failed to enable encryption:', e);
                btn.disabled = false;
                btn.textContent = 'Enable Encryption';
            }
        };
        encryptPanelContent.appendChild(btn);
        return;
    }

    const keyStr = localStorage.getItem('ohmycode_key_' + file.id);
    if (!keyStr || !app.encKey) {
        const msg = document.createElement('p');
        msg.className = 'encrypt-panel-msg';
        msg.textContent = 'This file is encrypted. Open the link that includes the key to manage sharing options.';
        encryptPanelContent.appendChild(msg);
        return;
    }

    const {editLink, roLink, editKey, roKey} = buildShareLinks();

    if (editKey) addLinkRow('Edit key:', editKey);
    addLinkRow('Edit link:', editLink);
    if (roLink) {
        if (roKey) {
            addLinkRow('Read-only key:', roKey);
        } else {
            const noRoKey = document.createElement('p');
            noRoKey.className = 'encrypt-panel-msg';
            noRoKey.textContent = 'Read-only key is not available in this browser. Regenerate to create a new one (old read-only links will stop working).';
            encryptPanelContent.appendChild(noRoKey);

            const regenBtn = document.createElement('button');
            regenBtn.textContent = 'Regenerate read-only key';
            regenBtn.className = 'encrypt-action-btn';
            regenBtn.onclick = async (e) => {
                e.stopPropagation();
                regenBtn.disabled = true;
                regenBtn.textContent = 'Generating…';
                try {
                    const {key: newRoKey, exported: newRoExported} = await generateKey();
                    app.roEncKey = newRoKey;
                    localStorage.setItem('ohmycode_rokey_' + file.id, newRoExported);
                    await actions.setContent(contentCodeMirror.getValue());
                    renderEncryptPanel();
                } catch(err) {
                    console.error('Failed to regenerate RO key:', err);
                    regenBtn.disabled = false;
                    regenBtn.textContent = 'Regenerate read-only key';
                }
            };
            encryptPanelContent.appendChild(regenBtn);
        }
        addLinkRow('Read-only link:', roLink);
    } else {
        const pending = document.createElement('p');
        pending.className = 'encrypt-panel-msg';
        pending.textContent = 'Generating read-only link…';
        encryptPanelContent.appendChild(pending);
    }

    const disableBtn = document.createElement('button');
    disableBtn.textContent = 'Disable Encryption';
    disableBtn.className = 'encrypt-action-btn encrypt-action-danger';
    disableBtn.onclick = () => {
        disableBtn.disabled = true;
        app.encKey = null;
        app.roEncKey = null;
        file.encrypted = false;
        localStorage.removeItem('ohmycode_key_' + file.id);
        localStorage.removeItem('ohmycode_rokey_' + file.id);
        actions.setEncrypted(false);
        actions.setContent(contentCodeMirror.getValue());
        hideEncryptPanel();
    };
    encryptPanelContent.appendChild(disableBtn);
};

const showNoKeyOverlay = () => {
    noKeyOverlay.style.display = 'flex';
    noKeyInput.value = '';
    noKeyError.textContent = '';
    noKeyInput.focus();
};

const hideNoKeyOverlay = () => {
    noKeyOverlay.style.display = 'none';
};

const submitKey = async () => {
    const keyStr = noKeyInput.value.trim();
    noKeyError.textContent = '';
    try {
        await importKey(keyStr);
        const storagePrefix = app.isROLink ? 'ohmycode_rokey_' : 'ohmycode_key_';
        localStorage.setItem(storagePrefix + file.id, keyStr);
        window.location.reload();
    } catch(e) {
        noKeyError.textContent = 'Invalid key. Make sure you pasted the full key.';
    }
};

document.getElementById('no-key-submit').addEventListener('click', submitKey);
noKeyInput.addEventListener('keydown', (e) => { if (e.key === 'Enter') submitKey(); });

encryptBtn.onclick = encryptPanelToggle;

document.addEventListener('click', (event) => {
    if (isEncryptPanelOpen &&
        !encryptPanel.contains(event.target) &&
        !encryptBtn.contains(event.target)) {
        hideEncryptPanel();
    }
});

onFileChange((f) => {
    if (isEncryptPanelOpen) {
        renderEncryptPanel();
    }
    if (f.encrypted && !app.encKey) {
        showNoKeyOverlay();
    } else {
        hideNoKeyOverlay();
    }
});
