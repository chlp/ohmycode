# CLAUDE.md — E2E tests (Playwright)

This file describes the E2E tests for OhMyCode and everything needed to run and extend them.

## Test structure

```
e2e/
├── editor.spec.js          — basic editor behavior (load, WS, CodeMirror)
├── static-files.spec.js    — cache busting: hash in URL, Cache-Control, immutable headers
├── collaboration.spec.js   — two clients on the same file stay in sync (lang, content, name)
├── file-rename.spec.js     — inline contenteditable file name editing
├── lang-switch.spec.js     — language dropdown and mode switching
├── ui-controls.spec.js     — header buttons, controls/status bar visibility, new file
├── code-execution.spec.js  — run_task_with_content → runner → result pane, clean result
├── lock-mode.spec.js       — locking a file makes it read-only for every viewer
├── readonly-share.spec.js  — encrypted read-only share links (generate, open, decrypt, no-edit)
├── playwright.config.js    — config (chromium, baseURL, timeout, retries)
└── package.json            — dependencies (@playwright/test)
```

Related test layers (not E2E):

| Layer | Location | Command |
|-------|----------|---------|
| Go unit + WS integration (api) | `api/internal/...` | `cd api && go test -race ./...` |
| Go unit (runner) | `runner/internal/...`, `runner/config/...` | `cd runner && go test -race ./...` |
| JS unit (Vitest) | `api/internal/api/client/js/test/` | `cd api/internal/api/client && npm test` |
| **E2E (Playwright)** | `e2e/` | see below |

## Before running

1. **The app must be running** at `http://localhost:52674/`

```bash
# Start API + MongoDB
docker compose -f api/docker/docker-compose.yml up -d --build --remove-orphans --force-recreate

# Start runners (Python, Go, Node.js, Java, PHP, PostgreSQL, MySQL, ...)
docker compose -f runner/docker/docker-compose.yml up -d --build --remove-orphans --force-recreate
```

Docker must be running. If Docker Desktop isn't up: `open -a Docker && sleep 10`.

2. **E2E dependencies** (one-time):

```bash
cd e2e
npm install
npx playwright install chromium
```

## Running tests

```bash
cd e2e

# All tests
npm test

# Verbose output
npx playwright test --reporter=list

# A single file
npx playwright test editor.spec.js

# Interactive UI mode (useful for debugging)
npm run test:ui

# Against a different environment
APP_URL=https://ohmycode.work npx playwright test
```

## What each spec checks

### `editor.spec.js`
- App returns 200, `#content` textarea is present in the DOM
- URL changes to `/[22-char base62 ID]` on load
- `body` transitions from `opacity:0` to `opacity:1` after DOMContentLoaded
- WebSocket connects to `/file`
- CodeMirror renders and is visible
- `#file-header` is visible
- Editor receives a file snapshot after the WS handshake

### `static-files.spec.js`
- `index.html` is served with `Cache-Control: no-cache` — passes in both modes
- The other 5 tests check hash-based versioning and **auto-skip** in dev mode (`use_dynamic_files=true`): `?v=HASH` in the URL, immutable Cache-Control, hash stability

All 6 tests should pass in prod mode (`use_dynamic_files=false`, embedded files).

### `collaboration.spec.js`
- Two browser contexts open the same file URL and see the same language, content, and file name

### `file-rename.spec.js`
- `#file-name` is visible and `contenteditable`
- Clicking it focuses it; typing a new name updates the visible text

### `lang-switch.spec.js`
- `#lang-select` is visible and contains the expected language options
- Selecting a language updates the select value across changes

### `ui-controls.spec.js`
- History/lock/encrypt/download/versions header buttons are visible
- `#controls-container` and `#status-bar` become visible after WS init
- "New file" navigates to a distinct new 22-char file URL

### `code-execution.spec.js`
- Selecting `python3` and clicking Run shows the program's output in the result pane
- Clean Result clears a previous run's output and hides the result pane
- Requires a runner (e.g. python3 container) to be online — the run button only enables once one is

### `lock-mode.spec.js`
- Clicking the lock button sets the status bar to "Locked" and makes the CodeMirror instance read-only
- A lock set by one client is reflected on another client viewing the same file

### `readonly-share.spec.js`
- Enabling encryption and copying the generated read-only link lets a second, unrelated browser context open it, decrypt the content client-side, and see it — but not edit it (lock button disabled, editor read-only)

## Adding new tests

New `.spec.js` files in `e2e/` are picked up automatically (see `testDir: '.'` in the config).

Typical template:

```js
import { test, expect } from '@playwright/test';

test('description', async ({ page }) => {
  await page.goto('/');
  await page.waitForLoadState('networkidle');
  // ...
});
```

**Important details:**

- `baseURL` defaults to `http://localhost:52674`, overridden via `APP_URL`
- Timeout: 30s per test, 1 retry on failure
- Chromium only (intentional — the app targets desktop Chrome)
- `retries: 1` — don't treat a test as broken on the first failure, especially WS-dependent ones

## Common failures and diagnostics

| Symptom | Cause | Fix |
|---------|-------|-----|
| `Connection refused` on run | App isn't running | Start docker compose (see above) |
| `ERR_MODULE_NOT_FOUND: @playwright/test` | Dependencies not installed | `npm install` in `e2e/` |
| `Cannot find package 'playwright'` | Running from the wrong directory | `cd e2e` before commands |
| A WS-dependent test fails on the first attempt | Race condition on connect | Normal — the retry will pass |
| `test-results/` shows up in git status | Missing `.gitignore` | Already present: `e2e/.gitignore` |

## For Claude

**Always check before running tests:**

```bash
curl -s -o /dev/null -w "%{http_code}" http://localhost:52674/
# Should return 200. If not, start docker compose first (see above).
```

**E2E working directory:** always `cd e2e` before `npx playwright`. Running from the repo root or from `api/` gives `ERR_MODULE_NOT_FOUND`.

**Docker daemon:** if you get `Cannot connect to the Docker daemon` — run `open -a Docker` and wait ~10s before retrying.

**Debug screenshots:** Playwright saves traces to `e2e/test-results/` on a failed retry. View a trace with:

```bash
npx playwright show-trace e2e/test-results/<test-name>/trace.zip
```

**Writing exploratory tests:** for UX checks, use `page.keyboard.insertText()` instead of `page.keyboard.type()` — `type()` emulates character-by-character input and can lose focus in CodeMirror before the WS connects (the editor is `readOnly: true` until the first snapshot arrives).

**Two CodeMirror instances on the page:** `.CodeMirror` resolves to 2 elements once a runner is connected (the second is the read-only result pane, themed `tomorrow-night-bright`). Always scope with a container selector, e.g. `#content-container .CodeMirror` or `#result-container .CodeMirror`, rather than relying on `.first()`.

**WS tests:** register `page.waitForEvent('websocket', ...)` **before** `page.goto()` — otherwise the connect event can fire before the listener is attached and the test will hang until timeout.
