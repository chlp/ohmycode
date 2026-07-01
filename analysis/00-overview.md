# OhMyCode — Project Analysis (2026-07-01)

Full-repo audit, committed as reference. Snapshot of code state at 2026-07-01 — will drift, re-verify before acting on it. Split into files:

- [01-api.md](01-api.md) — API backend (Go, MongoDB, WS)
- [02-runner.md](02-runner.md) — Runner manager + per-language Docker sandboxes
- [03-client.md](03-client.md) — Browser client (vanilla JS)
- [04-docs-tests-deploy.md](04-docs-tests-deploy.md) — docs, tests, CI/CD, dev workflow

## What is OhMyCode

Collaborative online code editor / mini-IDE, browser-based. Real-time multi-user editing (last-write-wins, no CRDT), 9 "languages" runnable in isolated Docker sandboxes (Python3, Node20, Go1.21, Java11, PHP8.2, MariaDB "mysql8", Postgres13, JSON via jq, Markdown via pandoc). MongoDB persistence, daily version snapshots, read-only share links with client-side AES-GCM encryption, file locking. Deployed to Fly.io (`ohmycode.fly.dev` / `ohmycode.work`).

## High-level architecture

```
Browser ──WS /file──► API (Go + MongoDB)
                          │
                     TaskStore ──WS /runner──► Runner Manager (Go)
                                                      │
                                          data/{lang}/requests/
                                                      │
                                          Language containers (Docker, network-isolated)
                                                      │
                                          data/{lang}/results/
                                                      │
                                     Runner Manager ──set_result──► API ──► Browser
```

Two Go modules (`api/`, `runner/`), zero framework vanilla-JS client, no bundler. Single author (Aleksei Rytikov), 509 commits, 2023-11 → present, worked in bursts.

## Fastest way to get oriented

1. `CLAUDE.md` (root) — build/run/test/deploy commands, architecture.
2. `api/PROTOCOL.md` — exact WS wire protocol for `/file` and `/runner`.
3. `analysis/04-docs-tests-deploy.md` — testing/CI/doc gaps, including risk items.
