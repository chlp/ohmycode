#!/bin/bash
set -uo pipefail
shopt -s nullglob

adduser --disabled-password restricted_user

mkdir -p tmp
chmod -R 744 tmp

while true; do
    for REQUEST_FILEPATH in requests/*; do
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
        su -c "timeout 5 sh -c 'cat \"${REQUEST_FILEPATH}\" | jq'" restricted_user 1>>"$OUT" 2>&1
        if [ $? -eq 124 ]; then
          echo -e "\n\n-------------------------\nTimeout reached, aborting\n-------------------------\n" >> "$OUT"
        fi
        rm -f -- "$REQUEST_FILEPATH"
        mv -- "$OUT" "results/$ID"
    done
    sleep 0.01
done
