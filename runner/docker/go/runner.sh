#!/bin/bash
set -uo pipefail
shopt -s nullglob

adduser --disabled-password restricted_user

cd /app
mkdir -p requests results

mkdir -p go
chmod -R 755 go
mkdir -p tmp
chmod -R 744 tmp

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
        SRC="go/$ID.go"
        mv -- "$REQUEST" "$SRC"
        chmod 755 -- "$SRC"
        su -c "cd /app && timeout 10 go run \"go/$ID.go\"" restricted_user 1>>"$OUT" 2>&1
        if [ $? -eq 124 ]; then
          echo -e "\n\n-------------------------\nTimeout reached, aborting\n-------------------------\n" >> "$OUT"
        fi
        rm -f -- "$SRC"
        mv -- "$OUT" "results/$ID"
    done
    if [ "$found" -eq 0 ]; then
        inotifywait -qq -t 300 -e create -e moved_to requests >/dev/null 2>&1 || true
    fi
done
