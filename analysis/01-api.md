# API backend audit (`api/`)

Go module. Single binary: serves SPA static files, two WS endpoints, in-memory file cache flushed to MongoDB, brokers execution tasks to Runner(s). No REST — everything except `/health` and static assets is WS with small JSON action protocol.

## Architecture

- Two trust domains over WS: `/file` (browsers), `/runner` (runner-manager processes, optional token auth).
- All mutation funnels through `model.File` methods: lock → mutate → fan-out signal to subscriber channels (buffered size 1, coalescing intentional).
- Readers never touch fields directly — always `Snapshot()` (lock-guarded copy), safe to serialize/pass around.

## Endpoints

- `GET /health` — pings Mongo (2s timeout); 503 `degraded` if ping fails or `persist_errors >= 5`. `{status, mongo, persist_errors, runners_online}`.
- `/file` WS (`filehandler.go`):
  - Client→server: `init, set_content, set_name, set_user_name, set_lang, set_runner, set_locked, set_encrypted, clean_result, run_task, run_task_with_content, get_versions, restore_version`.
  - Server→client: full `fileDTO` snapshot (throttled 1×/500ms/client, content only sent when changed), `versions` response, `open_file` (redirect after version-restore).
- `/runner` WS (`runnerhandler.go`):
  - Runner→server: `init` (runner_id, is_public, runner_token), `set_result` (lang, hash, result).
  - Server→runner: `[]Task` (content, lang, hash — file_id/runner_id hidden from wire), 2s heartbeat.
- Static files: embedded FS (prod) or disk (dev, `use_dynamic_files`). Build-wide SHA-256 hash computed at startup, rewrites `.js` imports + `?v=` for cache-busting. Path traversal guarded; 22-char base62 paths fall through to `index.html` (SPA routing) — any well-formed 22-char string serves the SPA even for nonexistent files.

## Data model

- **`model.File`** — central mutable aggregate. `sync.Mutex`, every mutator bumps `UpdatedAt` + fan-out signal. `SetContent` rejects if `IsLocked` or writer mismatch (`Writer` compared by **app_id**, not user_id — single-writer lock is per app-instance/session, not per-user). Content/result capped at `content_max_length_kb` (default 512KB). `IsUnused()` = no active subscribers AND idle >10min (guards against wall-clock jumps from host sleep).
- **`model.User`** — ephemeral presence `{ID, Name, TouchedAt}`, pruned after 5s inactivity.
- **`model.Task`** — in-memory only, not persisted. `Hash uint32` (custom rolling hash `ohMySimpleHash`) correlates runner `set_result` back to task without runner seeing file IDs.
- **`model.Runner`** — `{ID, IsPublic, CheckedAt}`; "online" = checked <5s ago.
- **`model.FileVersion`** — max 1/day/file (UTC calendar day), `Preview` = one-line diff summary vs last version, trimmed to latest 20.

## Store layer

- **FileStore** — `map[string]*model.File` + per-file mutex map (serializes per-file ops without blocking unrelated files). Mongo `ReplaceOne upsert` on persist. Tracks consecutive persist errors (surfaced via `/health`). `FlushAll()` on graceful shutdown.
- **TaskStore** — `map[string]*model.Task` keyed by FileId + pub/sub channels. Retry: unclaimed tasks re-offered after 30s (mutates `RunnerId`/`GivenToRunnerAt` as a side effect of the "get" call — a disconnecting runner claims a task for 30s before anyone else can retry it).
- **RunnerStore** — for public-runner files, `IsOnline` checks "is *any* public runner online", not the specific assigned `runner_id` — the field is essentially vestigial for online-check purposes.
- **VersionStore** — insert-if-new-day then trim to 20 newest. **No Mongo indices created anywhere in code** for `file_versions` (queried by `file_id` + sorted `created_at` — full scan risk at scale; may exist out-of-band on Atlas, unconfirmed).

## Background workers

- Cleanup: **every 1s** (root/api CLAUDE.md incorrectly imply 10min — that's the idle *threshold*, not the tick interval). Drops stale presence (>5s), releases write lock (idle >2s), force-timeouts stuck "waiting for result" (20s, message "Execution timed out after 20 seconds"), deletes unused files.
- Persist: every 30s, only dirty files.
- Runner-online sync: every 1s, recomputes `IsRunnerOnline`, no-op if unchanged.
- Shutdown: HTTP drain (10s) → `FlushAll()` sync → close Mongo. Clean shutdown guarantees no lost edits.

## Config

`api-conf.json` (falls back to `-example`) + env overrides: `OHMYCODE_MONGO_URI`, `OHMYCODE_MONGO_DBNAME`, `OHMYCODE_PORT`, `OHMYCODE_WS_ORIGINS`, `OHMYCODE_RUNNER_TOKEN`. `api-conf.json`/`api-conf.prod.json`/`fly-secrets.env` all gitignored, confirmed never committed.

## Recent fixes (git log signal)

Last ~20 commits mostly chase edge cases in the newest feature (RO share links + encryption + lock):
- `19d7e5a` / `de3d204` — `IsLocked` / `Encrypted`/`ROToken`/`ROContent` were never written into the Mongo persist doc → state reverted to unlocked/unencrypted after every restart until fixed.
- `3f17335` — bundled commit: added runner-token auth **and** markdown XSS sanitization together (see client findings — sanitization coverage is incomplete elsewhere).
- `8041faf` — separated edit-key vs RO-key (previously RO links may have shared the edit key — security-relevant redesign).
- `f48d46d` — Docker config-file-selection bug (wrong conf loaded in container).

## Risks (see also 05-risks-and-findings.md)

- **Runner token optional, defaults open** — anyone reaching `/runner` can register as a runner and receive tasks/content unless `OHMYCODE_RUNNER_TOKEN` is explicitly set in prod.
- No Mongo indices provisioned in code for version lookups.
- In-memory TaskStore/RunnerStore — restart drops in-flight tasks silently (client eventually sees 20s timeout, no explicit error).
- `IsRunnerOnline` staleness window (~6s) — `run_task` can be accepted right as a runner dies, surfaced only via the 20s timeout.
- No rate limit on most WS actions (only `run_task*`, 30/min/IP) — `set_content` unthrottled beyond 4MB/message + 512KB/file caps.
- Doc drift: cleanup interval documented as "10 min" (actually idle threshold; tick is 1s).
