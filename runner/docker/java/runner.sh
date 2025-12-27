#!/bin/bash
set -uo pipefail
shopt -s nullglob

adduser --disabled-password --gecos "" restricted_user

cd /app
mkdir -p requests results
chmod 755 requests results 2>/dev/null || true

mkdir -p java
chmod -R 755 java
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
        chmod 644 -- "$REQUEST" 2>/dev/null || true
        DIR="java/$ID"
        mkdir -p -- "$DIR"
        mv -- "$REQUEST" "$DIR/Main.java"
        chmod -R 755 -- "$DIR"
        timeout 10 javac "$DIR/Main.java" 1>"$OUT" 2>&1
        if [ $? -eq 0 ]; then
            (
                cd "$DIR"
                su -c "cd \"$PWD\" && timeout 10 java Main" restricted_user 1>>"../../$OUT" 2>&1
                if [ $? -eq 124 ]; then
                    echo -e "\n\n-------------------------\nTimeout reached, aborting\n-------------------------\n" >> "../../$OUT"
                fi
            )
            rm -rf -- "$DIR"
        else
            echo -e "\n\n-------------------------\nCompilation failed\n-------------------------\n" >> "$OUT"
        fi
        mv -- "$OUT" "results/$ID"
    done
    if [ "$found" -eq 0 ]; then
        inotifywait -qq -t 300 -e create -e moved_to requests >/dev/null 2>&1 || true
    fi
done
