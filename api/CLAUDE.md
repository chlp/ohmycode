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
go test ./internal/model
```

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

### Storage (internal/store/)

- **FileStore** — in-memory cache + MongoDB persistence, per-file mutex
- **TaskStore** — task queue with pub/sub for runner notification
- **RunnerStore** — runner registry with online status tracking

### Background Worker (internal/worker/)

- Cleanup of unused files (10 minutes)
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

Defined in `internal/model/lang.go`: go, java, json, markdown, mysql8, php82, postgres13

## Client Files

JavaScript modules in `internal/api/client/js/`:
- `app.js` — initialization
- `connect.js` — WebSocket management
- `editor.js` — CodeMirror integration
- `file.js` — file model
- `run.js` — code execution UI
