# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

Make the smallest possible code changes with the largest impact, thoroughly reasoned, anticipating side effects, and avoiding accidental influence.

## Project Overview

OhMyCode API is a real-time collaborative code editing and execution platform. Uses WebSocket for communication and MongoDB for data storage.

## Build & Run

```bash
# Configuration
cp api-conf-example.json api-conf.json
# Edit api-conf.json (MongoDB URI, port, etc.)

# Local build
go build -o api cmd/main.go
./api

# Docker (recommended)
docker compose -f docker/docker-compose.yml up --build
```

## Testing

```bash
# All Go tests (unit + WS integration), with race detector
go test -race ./...

# All tests including JS unit (from repo root)
../test.sh

# All tests including E2E (requires running app)
../test.sh --e2e
```

### Test layers

| Layer | Location | Command | Needs DB? |
|-------|----------|---------|-----------|
| Go model unit | `internal/model/*_test.go` | `go test -race ./...` | No |
| Go store unit | `internal/store/*_test.go` | `go test -race ./...` | No |
| Go API unit | `internal/api/health_test.go`, `headers_test.go`, `ratelimiter_test.go` | `go test -race ./...` | No |
| Go WS integration | `internal/api/ws_integration_test.go` | included above | No — uses `store.NewFileStoreInMemory()` |
| Go static files | `internal/api/staticfiles_test.go` | included above | No |
| JS unit (Vitest) | `internal/api/client/js/test/` | `cd internal/api/client && npm test` | No |
| E2E (Playwright) | `../../e2e/` | see root CLAUDE.md | App must run |

### Key test files

- **`internal/api/staticfiles_test.go`** — cache-busting hash, JS import patching, HTTP handler
- **`internal/api/health_test.go`** — `/health` endpoint: status, JSON body, runners count
- **`internal/api/headers_test.go`** — CORS headers, OPTIONS, cache-control per route
- **`internal/api/ratelimiter_test.go`** — run rate limiter, `clientIP`, WebSocket origin check
- **`internal/api/ws_integration_test.go`** — WS init→snapshot, set_content broadcast, set_lang, set_encrypted, RO token, invalid ID disconnect
- **`internal/model/file_test.go`** — `File` model: SetContent, SetLang, SetEncrypted, SetLocked, SetName, subscribe/unsubscribe, concurrency
- **`internal/store/filestore_test.go`** — in-memory store CRUD and concurrent access
- **`internal/store/runnerstore_test.go`** — RunnerStore: SetRunner, IsOnline, CountOnline, GetPublicRunner, TouchRunner
- **`internal/store/versionstore_test.go`** — `diffPreview`/`truncLine`/`splitTrimmedLines`: added/removed/reordered-line detection, 62-rune truncation, blank-line handling

## Architecture

### Core Components

- **cmd/main.go** — entry point
- **internal/api/apiservice.go** — HTTP server with two WebSocket endpoints:
  - `/file` — for clients (browsers)
  - `/runner` — for runner services (code execution)
- **internal/api/filehandler.go** — file operations handling (init, set_content, run_task)
- **internal/api/runnerhandler.go** — task distribution among runners
- **internal/api/wsclient.go** — WebSocket connection abstraction with ping/pong

### Data Flow

```
Browser → WS /file → API → FileStore (in-memory + MongoDB)
                        ↓
                   TaskStore → pub/sub → Runner (WS /runner)
                        ↓
                   Runner executes code → set_result → File.Result → Clients
```

### Models (internal/model/)

- **File** — document with content, metadata, users, execution result
- **Task** — execution task (content, language, hash, runner)
- **Runner** — code execution service (ID, public/private, online status)
- **FileVersion** — daily content snapshot (max 1/UTC day/file, trimmed to latest 20) for the version history panel

### Storage (internal/store/)

- **FileStore** — in-memory cache + MongoDB persistence, per-file mutex
- **TaskStore** — task queue with pub/sub for runner notification
- **RunnerStore** — runner registry with online status tracking
- **VersionStore** — MongoDB-backed version snapshots; queried by `file_id` + sorted by `created_at`. No index is created in code for `file_versions` — confirm one exists on the Mongo side at scale.

### Background Worker (internal/worker/)

- Cleanup sweep every 1 second — deletes files idle >10 minutes with no subscribers
- Persistence to MongoDB (every 30 seconds)
- Runner status synchronization (every second)

## Concurrency

- **File** uses `sync.Mutex` with snapshot pattern
- **Stores** use `sync.RWMutex`
- Per-file mutex in FileStore prevents lock contention
- Buffered channels for pub/sub

## Configuration (api-conf.json)

- `db.connectionString` — MongoDB URI
- `db.dbname` — database name
- `http_port` — server port (default: 52674)
- `serve_client_files` — serve client files
- `use_dynamic_files` — true for dev (from disk), false for prod (embedded)
- `ws_allowed_origins` — whitelist for WebSocket origins

## Supported Languages

Defined in `internal/model/lang.go`: go, java, json, markdown, mysql8, nodejs, php82, postgres13, python3

## Client Files

JavaScript modules in `internal/api/client/js/`:
- `main.js` — entry point (side-effect imports only)
- `app.js` — initialization, file open/navigate, UUID generation
- `connect.js` — WebSocket management, exponential backoff reconnect (recursive setTimeout: 1s → 2s → 4s → … → 30s)
- `editor.js` — CodeMirror integration, content sync, Cmd+S download, Cmd/Ctrl+Enter shortcut
- `file.js` — file state application, persistence triggers, IndexedDB save
- `file_name.js` — file name editing: saves immediately on blur, 5s auto-save while typing
- `lang.js` — language selector and CodeMirror mode switching; order: Python 3, Node.js, GoLang, Java, PHP, MySQL, PostgreSQL, JSON, Markdown
- `run.js` — code execution UI, result pane, status bar hint on hover when no runner connected
- `sidebar.js` — file history panel (with empty state message)
- `versions.js` — version history panel, restore-in-new-tab
- `hello_world.js` — "Code example" button: inserts language-specific starter code
- `open_file.js` — drag-and-drop file upload with binary detection and language auto-detection from extension
- `db.js` — IndexedDB cache for offline/fast load
- `status.js` — status bar: lock messages (Offline, Blocked), transient notifications (save time, run time), idle info (file size)
- `utils.js` — pure helpers (`ohMySimpleHash`)
- `encrypt.js` — AES-GCM 256 client-side encryption primitives (Web Crypto)
- `encrypt_ui.js` — encryption panel: enable/disable, edit/read-only key management, share-link generation, "no key" unlock overlay

CodeMirror modes bundled in `internal/api/client/codemirror/mode/`:
`clike`, `css`, `go`, `htmlmixed`, `javascript`, `markdown`, `php`, `python`, `sql`, `xml`

### Static file serving & cache busting

`internal/api/staticfiles.go` computes a SHA-256 hash of all embedded client files at
startup. Every relative JS import (`"./foo.js"`) is rewritten to `"./foo.js?v=HASH"` in
all served modules, and `?v=N` in `index.html` is replaced with `?v=HASH`. Versioned
assets get `Cache-Control: immutable`; `index.html` gets `no-cache`.
