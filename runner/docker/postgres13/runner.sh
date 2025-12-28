#!/bin/bash
set -uo pipefail
shopt -s nullglob

cd /app

mkdir -p requests results
chmod 755 requests results 2>/dev/null || true

mkdir -p tmp

export PGPASSWORD="${POSTGRES_PASSWORD:-password}"

while true; do
    found=0
    for REQUEST in requests/*; do
        found=1
        echo "$REQUEST"
        ID="$(basename -- "$REQUEST")"
        if ! [[ "$ID" =~ ^[0-9]+$ ]]; then
            echo "Invalid request id: $ID" >&2
            rm -f -- "$REQUEST"
            continue
        fi
        OUT="tmp/$ID"
        touch -- "$OUT"
        chmod 744 -- "$OUT"
        chmod 644 -- "$REQUEST" 2>/dev/null || true

        psql -U "${POSTGRES_USER:-user}" -d "${POSTGRES_DB:-mydatabase}" -c "CREATE DATABASE tmp_${ID};" 1>/dev/null 2>&1
        psql -U "${POSTGRES_USER:-user}" -d "tmp_${ID}" -c "CREATE USER tmp_user_${ID} WITH PASSWORD '${PGPASSWORD}';" 1>/dev/null 2>&1
        psql -U "${POSTGRES_USER:-user}" -d "tmp_${ID}" -c "GRANT ALL PRIVILEGES ON DATABASE tmp_${ID} TO tmp_user_${ID};" 1>/dev/null 2>&1

        timeout 5 psql -U "tmp_user_${ID}" -d "tmp_${ID}" -q -P format=wrapped -f "$REQUEST" 1>>"$OUT" 2>&1
        if [ $? -eq 124 ]; then
          echo -e "\n\n-------------------------\nTimeout reached, aborting\n-------------------------\n" >> "$OUT"
        fi

        psql -U "${POSTGRES_USER:-user}" -d "${POSTGRES_DB:-mydatabase}" -c "DROP DATABASE tmp_${ID};" 1>/dev/null 2>&1
        rm -f -- "$REQUEST"
        mv -- "$OUT" "results/$ID"
    done
    if [ "$found" -eq 0 ]; then
        inotifywait -qq -t 300 -e create -e moved_to requests >/dev/null 2>&1 || true
    fi
done
