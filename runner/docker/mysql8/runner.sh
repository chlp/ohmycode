#!/bin/bash
set -uo pipefail
shopt -s nullglob

cd /app

mkdir -p requests results
chmod 755 requests results 2>/dev/null || true

mkdir -p tmp

# MariaDB image uses MARIADB_ROOT_PASSWORD (we keep MYSQL_ROOT_PASSWORD as a backward-compatible alias).
export MYSQL_PWD="${MYSQL_ROOT_PASSWORD:-${MARIADB_ROOT_PASSWORD:-mysql8root}}"

while true; do
    mysql -u root -e "SHOW DATABASES;" 1>/dev/null 2>&1
    if [ $? -eq 0 ]; then
        echo "Starting mysql8 runner"
        break
    fi
    echo "Waiting for mysql8"
    sleep 2
done

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
        mysql -u root -e "CREATE DATABASE tmp_${ID};"
        timeout 5 mysql -u root "tmp_${ID}" --table < "$REQUEST" 1>>"$OUT" 2>&1
        if [ $? -eq 124 ]; then
          echo -e "\n\n-------------------------\nTimeout reached, aborting\n-------------------------\n" >> "$OUT"
        fi
        mysql -u root -e "DROP DATABASE tmp_${ID};"
        rm -f -- "$REQUEST"
        mv -- "$OUT" "results/$ID"
    done
    if [ "$found" -eq 0 ]; then
        inotifywait -qq -t 300 -e create -e moved_to requests >/dev/null 2>&1 || true
    fi
done
