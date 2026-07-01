# Client audit (`api/internal/api/client/`)

**Note**: real served client is `api/internal/api/client/` (embedded via `go:embed`, served dynamically from disk in dev). Top-level `client/public/` is empty/vestigial — ignore it.

## Stack

Vanilla JS, ES modules, no framework, no bundler. Third-party libs vendored (not npm): CodeMirror 5, `marked.min.js` + DOMPurify (`md/purify.min.js`). `package.json` in the client dir exists only for Vitest devDependency, not a real build. Cache-busting done server-side (see api findings).

## Module inventory (17 files under `js/`)

`main` (entry), `app` (file/app reactive state, routing, UUID), `connect` (WS lifecycle + protocol), `editor` (CodeMirror instance, lock/RO reflection, shortcuts), `file` (applies server snapshots, IndexedDB, "Saved" status), `file_name` (inline rename), `lang` (language registry/select), `run` (run button + result pane), `sidebar` (history panel), `versions` (version history/restore), `hello_world`, `open_file` (drag-drop upload), `db` (IndexedDB), `status` (status bar), `utils` (`ohMySimpleHash`), `encrypt` (AES-GCM primitives), `encrypt_ui` (encryption panel + RO share links).

`encrypt.js`/`encrypt_ui.js` are not mentioned in `api/CLAUDE.md`'s client file list — worth adding.

## Features

- CodeMirror 5 editor + separate read-only result-pane instance.
- Language dropdown: Python3, Node.js, Go, Java, PHP8.2, MySQL8, Postgres13, JSON, Markdown (edit/view).
- Locking: header toggle → `set_locked`; locked state disables editing, hides edit/view toggle, disables lang/hello-world/clean/run. Distinct from soft writer-lock (`file.writer_id`, per app-instance).
- RO share + client-side encryption: AES-GCM keys generated client-side, stored only in localStorage, never sent in plaintext. Separate edit-key vs RO-key. Server never sees plaintext for encrypted files.
- Autosave: no explicit save button — `contentSender` polls 500-1000ms, pushes `set_content` on hash mismatch. Server persists to Mongo every 30s.
- Shortcuts: Cmd/Ctrl+S download, Cmd/Ctrl+Enter run or markdown view/edit toggle, Esc toggles history panel.
- History (IndexedDB, flat list) distinct from Versions (server daily snapshots, restore opens in new tab).
- Drag-drop upload with binary-content heuristic + extension→language mapping.

## WS protocol / reconnect

Matches API's `/file` protocol (see 01-api.md). Reconnect: exponential backoff 1s→30s cap via recursive `setTimeout`, resets on success. Re-sends `init` on reconnect. Stale-file guard in `file.js`'s `applyFile` ignores snapshots whose `id` doesn't match currently-open file (comment cites past "UI reset loops" — this was a real historical bug class).

## State management

No framework — module-level singletons (`file`, `app` objects) with getter/setter properties as manual reactivity (e.g. `file.is_locked` setter drives UI updates). Cross-module wiring via small callback-array pub/sub (`onFileChange`, `onLangChange`, etc). IndexedDB mirrors file state; localStorage holds encryption keys/RO tokens/user identity.

## Test coverage

Only `js/test/utils.test.js` — 7 cases, all for `ohMySimpleHash` (the one pure function in the whole client). Everything else (DOM, WS protocol, reactivity, encryption, locking) is untested at unit level; pushed to Playwright E2E, which itself doesn't cover most of this (see 04-docs-tests-deploy.md).

## Recent fixes (git log)

Last 5 client commits are all small guard-clause fixes clustered around the lock/RO feature:
- `9b3ef20` — false "Saved" toast on page load (fired before `app.isOnline` was true).
- `696a625` — Cmd/Ctrl+Enter could bypass lock/RO view-switch guard (keyboard-only bypass of an existing button-level guard).
- `9b7149f` — null `file` during init crashed `updateEditViewButtons`.
- `7b58274` — introduced shared `updateEditViewButtons()`, added `is_locked` to IndexedDB (previously not cached → stale unlocked state after reload from cache).
- `3f17335` (security) — DOMPurify added around `marked.parse()` output in `file.js`'s server-content render path.
- `7387bec` (older) — fixed double-invoked callback in `connect.js`, fixed IndexedDB Promise wrapper (`tx.complete` doesn't exist on real `IDBTransaction` — was a silent no-op for awaiters).

## Risks (see also 05-risks-and-findings.md)

- **Markdown XSS sanitization is inconsistent.** `DOMPurify.sanitize` wraps `marked.parse()` only in `file.js`'s server-push render path. Three other call sites still do raw `innerHTML = marked.parse(...)` unsanitized: `editor.js:135` (every local keystroke in markdown mode — self-XSS risk if content came from elsewhere e.g. drag-drop), `open_file.js:67` (after drag-drop upload), `hello_world.js:212` (static content, low risk). Recommend patching all four for defense-in-depth.
- **Dead presence feature**: `file.users` still fully plumbed (server → state → IndexedDB) but `users.js` (the renderer) was deleted in an earlier "new left menu" commit; `#users-container` is an empty span. Cleanup or restore.
- **Collaboration is last-write-wins, not CRDT/OT.** Writer exclusivity (`file.writer_id`) is advisory/soft; concurrent edits within one round-trip window can clobber each other. No remote cursor/selection indicators (removed along with presence UI).
- **Duplicated "can edit" logic**: `app.isROLink || file.is_locked` checks are independently repeated across `editor.js`, `lang.js`, `run.js`, `hello_world.js`, `encrypt_ui.js` instead of one shared derived flag — the exact pattern that produced the last 5 bugfix commits (missed guard in one of N places). Recommend consolidating.
- **RO-key regeneration silently breaks existing share links** with no warning to holders of old links (edge case in `encrypt_ui.js`, e.g. user clears localStorage on new device).
- Two `// todo:` markers in shipped code: markdown-only hardcoding of view/edit toggle (`editor.js`, `lang.js`) — mechanism is generic but restricted to one language.
