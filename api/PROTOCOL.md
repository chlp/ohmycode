# WebSocket Protocol

OhMyCode API exposes two WebSocket endpoints:

- `/file` — for browser clients
- `/runner` — for code-execution runner services

All messages are JSON text frames.

---

## `/file` — Client ↔ API

### Client → API

| Action | Required fields | Description |
|--------|----------------|-------------|
| `init` | `file_id`, `user_id`, `app_id`, `file_name`, `lang`, `content` | Open or create a file. Must be the first message sent. |
| `set_content` | `content` | Update file content. |
| `set_name` | `file_name` | Rename the file. |
| `set_lang` | `lang` | Change the language. |
| `set_runner` | `runner_id` | Assign a private runner by ID. |
| `set_user_name` | `user_name` | Set the display name for the current user. |
| `run_task` | _(none)_ | Execute the current file content. Ignored if execution is already in progress. |
| `run_task_with_content` | `content` | Atomically update content and execute. |
| `clean_result` | _(none)_ | Clear the execution result. |
| `get_versions` | _(none)_ | Request the version history list. |
| `restore_version` | `version_id` | Restore a version into a new file. |

**Example — init:**
```json
{"action":"init","file_id":"abc123...","user_id":"uuid...","app_id":"uuid...","file_name":"main.py","lang":"python3","content":""}
```

**Example — run_task_with_content:**
```json
{"action":"run_task_with_content","content":"print('hello')"}
```

### API → Client

The server pushes a file snapshot whenever the file state changes. Snapshots are throttled to at most one per 500 ms per subscriber.

**File snapshot** (sent on init and on every state change):
```json
{
  "id": "abc123...",
  "name": "main.py",
  "lang": "python3",
  "content": "print('hello')",
  "content_updated_at": "2026-05-31T12:00:00Z",
  "result": "hello\n",
  "writer_id": "uuid...",
  "runner": "uuid...",
  "users": [{"id":"uuid...","name":"Alice","touched_at":"..."}],
  "updated_at": "2026-05-31T12:00:01Z",
  "persisted": true,
  "is_waiting_for_result": false,
  "is_runner_online": true
}
```

> `content` is omitted in snapshots that contain no content change (only metadata changed).

**Versions list** (response to `get_versions`):
```json
{"action":"versions","versions":[{"id":"uuid...","name":"main.py","lang":"python3","created_at":"..."}]}
```

**Open file** (response to `restore_version`):
```json
{"action":"open_file","file_id":"newFileId..."}
```

---

## `/runner` — Runner ↔ API

### Runner → API

| Action | Required fields | Description |
|--------|----------------|-------------|
| `init` | `runner_id`, `is_public` | Register this runner. Must be the first message. |
| `set_result` | `lang`, `hash`, `result` | Return execution output for a task. |

**Example — init:**
```json
{"action":"init","runner_id":"uuid...","is_public":false}
```

**Example — set_result:**
```json
{"action":"set_result","lang":"python3","hash":1234567890,"result":"hello\n"}
```

### API → Runner

The server pushes task lists whenever a new task is available or on a 2-second heartbeat.

**Task list:**
```json
[
  {
    "file_id": "abc123...",
    "content": "print('hello')",
    "lang": "python3",
    "hash": 1234567890,
    "runner_id": "uuid...",
    "is_public": false
  }
]
```

---

## Execution flow

```
Client                 API                 Runner
  |                     |                    |
  |── run_task_with ──► |                    |
  |   _content          |── [task list] ───► |
  |                     |                    |── executes code
  |◄── snapshot ────────|◄─── set_result ───|
  | (is_waiting=true)   |                    |
  |◄── snapshot ────────|
  | (is_waiting=false,
  |  result="hello\n")
```

If the runner does not return a result within **20 seconds**, the API sets `result` to `"Execution timed out after 20 seconds"` and clears `is_waiting_for_result`.
