#!/bin/bash

exec docker-entrypoint.sh "$@" &

export MYSQL_PWD=

while true; do
    mysql -e "SHOW DATABASES;" 1>/dev/null 2>&1
    if [ $? -eq 0 ]; then
        echo "Starting mysql8 runner"
        break
    fi
    echo "Waiting for mysql8"
    sleep 2
done

mkdir -p tmp
while [ True ]; do
    if [ ! -z "$(ls tmp)" ]; then
      rm tmp/*
    fi
    if [ ! -z "$(ls requests)" ]; then
        for REQUEST in requests/*; do
            echo $REQUEST
            ID=$(basename $REQUEST)
            mysql -e "CREATE DATABASE tmp_$ID;"
            timeout 5 mysql tmp_$ID --table < $REQUEST 1>tmp/$ID 2>&1
            if [ $? -eq 124 ]; then
              echo -e "\n\n-------------------------\nTimeout reached, aborting\n-------------------------\n" >> tmp/$ID
            fi
            mysql -e "DROP DATABASE tmp_$ID;"
            mv tmp/$ID results/$ID
            rm $REQUEST
        done
    fi
    sleep 0.01
done
