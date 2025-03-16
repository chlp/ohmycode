const copyToClipboard = (text) => {
    if (navigator.clipboard && window.isSecureContext) {
        return navigator.clipboard.writeText(text).then(() => {
            console.log("Text copied to clipboard");
        }).catch(err => {
            console.error("Failed to copy: ", err);
        });
    } else {
        let textArea = document.createElement("textarea");
        textArea.value = text;

        textArea.style.position = "fixed";
        textArea.style.left = "-999999px";
        document.body.appendChild(textArea);
        textArea.focus();
        textArea.select();

        try {
            document.execCommand('copy');
            console.log("Text copied to clipboard");
        } catch (err) {
            console.error("Failed to copy: ", err);
        } finally {
            document.body.removeChild(textArea);
        }
    }
};
document.getElementById('sidebar-copy-url').onclick = async () => {
    await copyToClipboard(window.location.href);
};
