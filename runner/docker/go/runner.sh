#!/bin/bash
set -uo pipefail
shopt -s nullglob

adduser --disabled-password restricted_user

mkdir -p go
chmod -R 755 go
mkdir -p tmp
chmod -R 744 tmp

while true; do
    for REQUEST in requests/*; do
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
        su -c "timeout 10 go run \"go/$ID.go\"" restricted_user 1>>"$OUT" 2>&1
        if [ $? -eq 124 ]; then
          echo -e "\n\n-------------------------\nTimeout reached, aborting\n-------------------------\n" >> "$OUT"
        fi
        rm -f -- "$SRC"
        mv -- "$OUT" "results/$ID"
    done
    sleep 0.01
done
