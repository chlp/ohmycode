# Runner subsystem audit (`runner/`)

Stateless Go bridge: API ⟷ per-language Docker sandboxes. Never executes code itself — shuttles files.

## Architecture

- `internal/api/client.go` — WS client to API `/runner`, owns reconnect + outbound task queue.
- `internal/worker/worker.go` — one `TaskDistributor` + one `ResultProcessor` goroutine per configured language.
- File-based handoff to containers via bind-mounted `data/{lang}/{requests,results}/`.

## API connection

- `init` sends `{action, runner_id, is_public}` — **no `runner_token` field in the Go client struct**, even though the API supports an optional shared-token check. If a deployment ever sets `OHMYCODE_RUNNER_TOKEN`, this runner build cannot authenticate. Confirm whether intentional/unused.
- Reconnect: exponential backoff capped 30s + 0-3s jitter, resets on success. 60s pong-wait, 54s ping ticker.
- Read loop decodes `[]*Task` into local queue; `SetResult()` writes `set_result` JSON, closes socket on write failure to force reconnect.

## Task distribution protocol

- Requests: `TaskDistributor` writes `task.Content` to `data/{lang}/requests/<hash>.tmp`, `chmod 644`, then `os.Rename` → `<hash>` (atomic; container sees no partial writes). Filename = task's `Hash` (uint32 decimal), generated API-side, round-trips as correlation ID.
- Results: `ResultProcessor` polls `data/{lang}/results/`, parses filename as hash, reads content, calls `SetResult`, deletes file. Bad filename/read error → deletes file, substitutes `"something wrong with result"`.
- Both distributor/processor: adaptive polling 100ms → backs off ×10 to 1s when idle/erroring, resets on work.
- Inside containers: `while true` over `requests/*`, filename validated `^[0-9]+$` (rejects path traversal), executes, writes result, `mv`s into `results/`. Blocks on `inotifywait -t 300 -e create -e moved_to` when empty (matches manager's `os.Rename` → `moved_to` event).

## Per-language containers

| Lang | Base image | Exec | Timeout | Isolation note |
|---|---|---|---|---|
| go | golang:1.21-alpine | `go run` | 10s | compile+run combined, one timeout |
| java | eclipse-temurin:11-jdk (not slim) | `javac`+`java`, hardcoded class `Main` | 10s×2 | compile/run separate, distinct failure banner |
| json | alpine + jq | `jq` (formatter, not really a "language") | 5s | |
| markdown | alpine + pandoc | MD→RST conversion | 5s | |
| **mysql8** | **mariadb:10.11** (misnamed — it's MariaDB not MySQL8) | SQL run as **root** in ephemeral `tmp_<id>` DB | 5s | **no per-request user — root execution** |
| postgres13 | postgres:13.15 | SQL run as scoped **`tmp_user_<id>`** in ephemeral DB, DB+user dropped after | 5s | proper per-request user isolation |
| php82 | php:8.2-cli-alpine | `php` as `restricted_user` | 5s | |
| python3 | python:3.12-alpine | `python3` as `restricted_user` | 10s | |
| nodejs | node:20-alpine | `node` as `restricted_user` | 10s | |

Common pattern: unprivileged `restricted_user` runs the actual submitted code via `su`; outer polling loop runs as container default user (root, no `USER` directive).

## Isolation layers

1. **Network**: each language container on its own `internal: true` bridge network — no egress, no lateral movement between sandboxes. Strongest control present.
2. **Filesystem**: tmpfs for working dirs/DB data (wiped on restart); requests/results bind-mounts scoped to one manager↔one-language pair.
3. **Process user**: `restricted_user` for most languages (mysql8/postgres13 run their engine + user SQL as DB-superuser context — see gap below).
4. **Resource caps**: `deploy.resources.limits` (CPU/mem) declared in compose — **this is a Swarm-mode construct; under plain `docker compose up` enforcement is version/engine-dependent, not guaranteed.** Verify with `docker stats` under load rather than assume.
5. **Hard wall-clock `timeout`** per execution (5-10s) — no separate CPU/PID/fd ulimits, no seccomp/AppArmor customization.
6. Filename validation (`^[0-9]+$`) guards path traversal both runner-side and container-side.

## Config

`runner-conf.json` (bootstrapped from `-example` on first run, ID persisted via crypto/rand `GenUuid()`). Env overrides: `OHMYCODE_API_URL`, `OHMYCODE_RUNNER_ID` only — no override for `is_public`/`languages`/token. `runner-conf-example.json` has dead `api2`-`api5` keys (unused leftovers).

## Risks (see also 05-risks-and-findings.md)

- **mysql8 runs user-submitted SQL as root** inside its (network-isolated) container — no scoped user like postgres13's `tmp_user_<id>`. Asymmetric isolation between the two DB "languages"; recommend aligning mysql8 to a scoped temp user.
- **Runner-side token support missing** — protocol gap vs. API's optional auth check (see api findings).
- **`is_public` self-asserted, no server-side vetting** beyond the (currently-unimplementable-from-this-client) token — combined with the token gap, a rogue "public" runner could receive/execute other users' code if `OHMYCODE_RUNNER_TOKEN` isn't actually enforced in prod. **Worth confirming directly against the deployed API config.**
- **Zero test files anywhere under `runner/`** — CI only `go build`s it, never `go test`s (there's nothing to test). Reconnect/backoff and atomic file-rename distribution logic is unverified by automation.
- `deploy.resources.limits` reliability unconfirmed under actual compose engine version in use.
- Non-slim images (java `-jdk` not `-jre`, postgres/mariadb full Debian) — larger patch surface than needed for pure execution.
- Silent `|| true` swallowing on DB create/drop in mysql8/postgres13 `runner.sh` — a failed `DROP DATABASE` leaves orphaned `tmp_*` DBs unnoticed until container restart (tmpfs-backed, so bounded, but not logged).
