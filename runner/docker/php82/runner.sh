#!/bin/bash
set -uo pipefail
shopt -s nullglob

adduser --disabled-password restricted_user

cd /app
mkdir -p requests results tmp
chmod 755 requests results 2>/dev/null || true

while true; do
    found=0
    for REQUEST_FILEPATH in requests/*; do
        found=1
        echo "$REQUEST_FILEPATH"
        ID="$(basename -- "$REQUEST_FILEPATH")"
        if ! [[ "$ID" =~ ^[0-9]+$ ]]; then
            echo "Invalid request id: $ID" >&2
            rm -f -- "$REQUEST_FILEPATH"
            continue
        fi
        OUT="tmp/$ID"
        touch -- "$OUT"
        chmod 744 -- "$OUT"
        # Ensure restricted_user can read the request even if it was created with restrictive permissions on the host/volume.
        chmod 644 -- "$REQUEST_FILEPATH" 2>/dev/null || true
        su -c "cd /app && timeout 5 php \"${REQUEST_FILEPATH}\"" restricted_user 1>>"$OUT" 2>&1
        if [ $? -eq 124 ]; then
          echo -e "\n\n-------------------------\nTimeout reached, aborting\n-------------------------\n" >> "$OUT"
        fi
        rm -f -- "$REQUEST_FILEPATH"
        mv -- "$OUT" "results/$ID"
    done
    if [ "$found" -eq 0 ]; then
        # Idle: block on filesystem events instead of busy polling.
        inotifywait -qq -t 300 -e create -e moved_to requests >/dev/null 2>&1 || true
    fi
done
