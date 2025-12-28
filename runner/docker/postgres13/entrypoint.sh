#!/bin/bash
set -euo pipefail

# Start the request runner in background.
/app/runner.sh &

# Hand off to the original Postgres entrypoint (initdb, then exec postgres).
exec /usr/local/bin/docker-entrypoint.sh "$@"


