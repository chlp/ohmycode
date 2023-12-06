#!/bin/bash

exec docker-entrypoint.sh "$@" &
sleep 15 # waiting for mysql to start
export MYSQL_PWD=root_password

while [ True ]; do
    if [ ! -z "$(ls requests)" ]; then
        for REQUEST in requests/*; do
            ID=$(basename $REQUEST)
            mysql -e "CREATE DATABASE $ID;"
            mysql $ID --table < $REQUEST 1>results/$ID 2>&1
            mysql -e "DROP DATABASE $ID;"
            rm $REQUEST
        done
    fi
    sleep 1
done
