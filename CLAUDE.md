# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run

Run both services together (recommended):
```bash
docker compose -f api/docker/docker-compose.yml up -d --build --remove-orphans --force-recreate && \
docker compose -f runner/docker/docker-compose.yml up -d --build --remove-orphans --force-recreate
```
App available at http://localhost:52674/

Run services separately — see `api/CLAUDE.md` and `runner/CLAUDE.md`.

## Testing

```bash
# All tests (Go unit + WS integration + JS unit) — no external deps required
./test.sh

# Include E2E — requires app running at localhost:52674
./test.sh --e2e
APP_URL=https://ohmycode.work ./test.sh --e2e
```

Individual commands:

```bash
# Go — unit + WS integration (no external deps required)
cd api && go test -race ./...

# JS unit (Vitest, pure functions)
cd api/internal/api/client && npm test

# E2E (Playwright, requires running app at localhost:52674)
cd e2e && npm install && npx playwright install chromium && npx playwright test
# Override base URL:  APP_URL=https://ohmycode.work npx playwright test
```

### Test coverage

| Layer | Files | Needs running app? |
|-------|-------|--------------------|
| Go model unit | `internal/model/*_test.go` | No |
| Go store unit | `internal/store/*_test.go` | No |
| Go API unit (health, headers, rate limiter, origin) | `internal/api/*_test.go` | No |
| Go WS integration | `internal/api/ws_integration_test.go` | No |
| Go static files | `internal/api/staticfiles_test.go` | No |
| JS unit (Vitest) | `internal/api/client/js/test/` | No |
| E2E editor | `e2e/editor.spec.js` | Yes |
| E2E static files | `e2e/static-files.spec.js` | Yes |
| E2E lang switch | `e2e/lang-switch.spec.js` | Yes |
| E2E file rename | `e2e/file-rename.spec.js` | Yes |
| E2E collaboration | `e2e/collaboration.spec.js` | Yes |
| E2E UI controls | `e2e/ui-controls.spec.js` | Yes |

## Deploy

API deploys to Fly.io (`ohmycode` app → https://ohmycode.fly.dev/):

```bash
cd api && ./deploy.sh
```

Script sources `api/fly-secrets.env` then runs `flyctl deploy --remote-only` (remote build, no local Docker needed).

## Architecture

Two independent Go modules communicate over WebSocket:

```
Browser ──WS /file──► API (Go + MongoDB)
                          │
                     TaskStore ──WS /runner──► Runner Manager (Go)
                                                      │
                                          data/{lang}/requests/
                                                      │
                                          Language containers (Docker)
                                                      │
                                          data/{lang}/results/
                                                      │
                                     Runner Manager ──set_result──► API ──► Browser
```

**`api/`** — HTTP + WebSocket server, in-memory file cache backed by MongoDB.
- `/file` WS endpoint for browser clients
- `/runner` WS endpoint for runner services
- Files live in memory with per-file mutex; persisted to MongoDB every 30s
- Version snapshots taken once per day per file
- Background worker: cleanup (10 min), persist (30 s), runner status sync (1 s)

**`runner/`** — Stateless bridge between API and language containers.
- Connects to API via WebSocket with auto-reconnect
- Distributes tasks by writing files to `data/{lang}/requests/`
- Language containers read requests and write results to `data/{lang}/results/`
- Each language runs in an isolated Docker container with no external network access

### Key concurrency patterns in `api/`

- `File` uses `sync.Mutex` + snapshot pattern (`File.Snapshot()`) — never read fields directly, always snapshot under lock
- `FileStore` uses per-file mutex (`fileLocks`) to serialize GetOrCreate without blocking unrelated files
- `File.subs` is a fan-out map of per-subscriber buffered channels for WS push notifications

### IDs

Files use 22-character base62 IDs (not UUIDs). Users/runners use standard UUIDs. See `api/pkg/util/uuid.go`.

## Sub-project guidance

- `api/CLAUDE.md` — detailed API config, data flow, and client JS structure
- `runner/CLAUDE.md` — runner config, Docker setup per language

## Related repositories

- `../ohmycode-private` — private repo with product planning, ideas, and business notes (not public)
