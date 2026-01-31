# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

Make the smallest possible code changes with the largest impact, thoroughly reasoned, anticipating side effects, and avoiding accidental influence.

## Build & Run Commands

```bash
# Build the Go binary
go build -o app cmd/main.go

# Run the application
./app

# Run with Docker Compose (recommended for full setup)
cd docker && docker-compose up --build

# Build and run only specific services
cd docker && docker-compose up --build ohmycode-runner-manager ohmycode-runner-go
```

## Architecture

OhMyCode Runner is a distributed code execution service that:
1. Connects to a central API server via WebSocket
2. Receives code execution tasks
3. Distributes tasks to language-specific Docker containers via filesystem
4. Collects results and sends them back to the API

### Core Components

- **cmd/main.go** — Entry point; initializes API client and worker
- **config/config.go** — Loads `runner-conf.json` (auto-generated from `runner-conf-example.json` on first run)
- **internal/api/client.go** — WebSocket client with auto-reconnection; manages task queue
- **internal/worker/worker.go** — Orchestrates TaskDistributor and ResultProcessor goroutines
- **internal/worker/taskdistributor.go** — Writes incoming tasks to `data/{lang}/requests/` as files
- **internal/worker/resultprocessor.go** — Reads results from `data/{lang}/results/` and sends to API

### Data Flow

```
API Server ←→ WebSocket ←→ Runner Manager
                              ↓
                    data/{lang}/requests/  (task files)
                              ↓
              Language Containers (go, java, php82, etc.)
                              ↓
                    data/{lang}/results/   (result files)
                              ↓
                    Runner Manager → API Server
```

### Configuration

`runner-conf.json` fields:
- `id` — 32-char UUID (auto-generated if missing)
- `is_public` — Whether runner is publicly available
- `languages` — Array of supported languages (must match docker-compose services)
- `api` — WebSocket URL of the API server

### Docker Setup

Each language runs in an isolated container with:
- Internal-only network (no external access)
- Mounted volumes for requests/results at `/app/requests` and `/app/results`
- tmpfs for temporary execution files

Supported languages: go, java, json, markdown, mysql8, php82, postgres13
