const generateKey = async () => {
    const key = await crypto.subtle.generateKey(
        {name: 'AES-GCM', length: 256},
        true,
        ['encrypt', 'decrypt']
    );
    const exported = await exportKey(key);
    return {key, exported};
};

const exportKey = async (key) => {
    const raw = await crypto.subtle.exportKey('raw', key);
    return btoa(String.fromCharCode(...new Uint8Array(raw)))
        .replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '');
};

const importKey = async (b64) => {
    const standard = b64.replace(/-/g, '+').replace(/_/g, '/');
    const padded = standard + '='.repeat((4 - standard.length % 4) % 4);
    const bytes = Uint8Array.from(atob(padded), c => c.charCodeAt(0));
    return crypto.subtle.importKey('raw', bytes, {name: 'AES-GCM'}, true, ['encrypt', 'decrypt']);
};

const encryptText = async (key, plaintext) => {
    const iv = crypto.getRandomValues(new Uint8Array(12));
    const encoded = new TextEncoder().encode(plaintext);
    const ciphertext = await crypto.subtle.encrypt({name: 'AES-GCM', iv}, key, encoded);
    const combined = new Uint8Array(12 + ciphertext.byteLength);
    combined.set(iv);
    combined.set(new Uint8Array(ciphertext), 12);
    return btoa(String.fromCharCode(...combined));
};

const decryptText = async (key, b64) => {
    const combined = Uint8Array.from(atob(b64), c => c.charCodeAt(0));
    const iv = combined.slice(0, 12);
    const ciphertext = combined.slice(12);
    const decrypted = await crypto.subtle.decrypt({name: 'AES-GCM', iv}, key, ciphertext);
    return new TextDecoder().decode(decrypted);
};

export {generateKey, exportKey, importKey, encryptText, decryptText};
