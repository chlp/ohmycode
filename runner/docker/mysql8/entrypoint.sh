#!/bin/bash
set -euo pipefail

# Start the request runner in background; it will wait until MySQL is ready.
/app/runner.sh &

# Hand off to the original MySQL entrypoint (initdb, then exec mysqld).
exec /usr/local/bin/docker-entrypoint.sh "$@"


