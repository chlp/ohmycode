# Docs, testing, CI/CD, dev workflow audit

## Documentation inventory

- `README.md` — marketing-quality pitch, screenshot, Medium article link, quick-start.
- `CLAUDE.md` (root) — build/run/test/deploy, architecture, concurrency, ID scheme. Updated in lockstep with features.
- `api/CLAUDE.md` — most substantial doc: components, data flow, models, storage, worker, config keys, languages, client JS module list, cache-busting.
- `runner/CLAUDE.md` — architecture, config, Docker isolation notes, languages.
- `api/PROTOCOL.md` — precise WS wire-protocol spec for `/file` and `/runner`, message tables, sequence diagram. High quality, hand-maintained.
- `e2e/CLAUDE.md` — **stale + written in Russian** (inconsistent with every other doc, which is English). Documents only `editor.spec.js` + `static-files.spec.js` in detail; silent on 4 other spec files (`collaboration`, `file-rename`, `lang-switch`, `ui-controls`). Contains a leftover hardcoded path (`/Users/amini/Documents/github/ohmycode`) from a different contributor's machine. Does contain genuinely useful Playwright gotchas (WS listener-before-navigate ordering, `.first()` disambiguation, `insertText` vs `type()` timing) — worth preserving those, just needs cleanup/translation/sync.
- `LICENSE` — MIT, Aleksei Rytikov 2023.
- `AGENTS.md` — symlink to `CLAUDE.md` for cross-tool compatibility.
- No CONTRIBUTING.md, CHANGELOG.md, or diagrams beyond ASCII art in CLAUDE.md.

**Assessment**: docs are unusually strong for project size — CLAUDE.md files visibly kept in sync with code (docs commits interleaved with feature commits). Main gap is `e2e/CLAUDE.md`.

## Testing layers

| Layer | Location | Needs infra? |
|---|---|---|
| Go model/store/api unit | `api/internal/{model,store,api}/*_test.go` (8 files) | No |
| Go WS integration | `api/internal/api/ws_integration_test.go` | No (in-memory store) |
| JS unit (Vitest) | `client/js/test/utils.test.js` | No |
| E2E (Playwright) | `e2e/*.spec.js` (6 files) | Yes, live app :52674 |
| **`runner/` tests** | **none exist** | — |

E2E spec coverage (read in full):
- `editor.spec.js` (7) — page load, redirect to file ID URL, WS connect, CodeMirror render, snapshot round-trip.
- `static-files.spec.js` (6) — cache headers; 5 of 6 auto-skip in dev mode (only run against embedded/prod build).
- `collaboration.spec.js` (3) — two browser contexts, content/lang/name sync within ~2s.
- `file-rename.spec.js` (4) — inline contenteditable rename flow.
- `lang-switch.spec.js` (4) — dropdown options + value switching.
- `ui-controls.spec.js` (7) — header buttons visible, controls appear post-WS-init, new-file navigation.

**Coverage gaps** — notably, these are exactly the newest/most-hardened features:
- No E2E for **code execution** (`run_task` → runner → result pane) despite being a core feature.
- No E2E for **lock mode / RO share links / encryption** despite ~10 recent commits dedicated to exactly this and a real historical bugfix tail (`de3d204`, `19d7e5a`, `9b3ef20`, `696a625`, `9b7149f`).
- No E2E for **version history/restore**, **drag-drop upload**, **runner offline/reconnect**, **rate limiting**.
- `runner/` has zero automated tests (CI only `go build`s it).
- No visible regression test specifically for the markdown XSS sanitization fix (`3f17335`).

## CI (`.github/workflows/test.yml`)

Push/PR to `main`, 3 parallel jobs: `go-api` (`go test -race`), `go-runner` (build-only), `js-unit` (Vitest). No E2E in CI (reasonable — needs Docker+Mongo+containers). No linting/static-analysis job (no golangci-lint/ESLint/`go vet`). No Dependabot/Renovate.

## Deployment

- Fly.io, app `ohmycode` → `ohmycode.fly.dev` / `ohmycode.work`.
- `.github/workflows/deploy.yml` — **manual `workflow_dispatch` only** (deliberately changed away from auto-deploy-on-push per `c3b9d9f`), runs `flyctl deploy --remote-only`.
- Local: `api/deploy.sh` sources gitignored `api/fly-secrets.env`, then `flyctl deploy --remote-only`.
- No staging environment, no rollback automation (relies on Fly.io release history), no version tags — hard to correlate a live deploy with a commit after the fact.
- Runner has no CI/CD path — self-hosted/manual by design ("bring your own private runner").

## Local dev workflow

Two separate `docker compose` invocations (api + runner), no root compose file/Makefile tying them together. Runner compose brings up **10 containers** (manager + 9 languages) by default — heavier than typical for a side project; can scope to specific services manually. Config bootstrap is manual copy of `*-example` files, no setup script. MongoDB/DB default passwords are plaintext in committed compose files (local-dev-only, behind `internal: true` networks — low real risk but worth knowing).

## Recent trajectory (git log, 40 commits reviewed / 509 total)

Coherent feature arc, not scattered: RO/encryption/lock feature (majority of recent commits) → persistence-correctness fixes → small UX polish. Earlier: deploy pipeline stood up then hardened (manual-trigger-only), then operational hardening (health checks, graceful shutdown, structured logging, rate limiting) *before* the higher-risk sharing feature — sensible ordering. Single author, worked in intense bursts (e.g. 24 commits in one day) separated by multi-week/month gaps.

## Secrets hygiene — checked, fine

`fly-secrets.env`, `api-conf.prod.json`, `api-conf.json`, `runner-conf.json`'s real ID: confirmed all either gitignored-and-never-committed, or (runner-conf.json) low-severity if committed since it's just an ID meant to be known to the API anyway. Only plaintext credential-shaped content in git is dev-only default DB passwords in compose files.

## Summary gaps to flag

1. Biggest recent feature areas (RO/encryption/lock, code execution) have **no E2E coverage** — exactly where bugfix churn has been highest.
2. `runner/` module untested.
3. `e2e/CLAUDE.md` stale + wrong language + leaked path.
4. No CI linting.
5. No dependency-update automation (relevant given this is a code-execution product).
6. No changelog/release versioning.
7. Solo dev, no code review — acceptable for scope, but no second pair of eyes on sandbox/security code.
