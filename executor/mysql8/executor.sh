#!/bin/bash

exec docker-entrypoint.sh "$@" &
sleep 15 # waiting for mysql to start
export MYSQL_PWD=root_password

mkdir tmp
while [ True ]; do
    if [ ! -z "$(ls tmp)" ]; then
      rm tmp/*
    fi
    if [ ! -z "$(ls requests)" ]; then
        for REQUEST in requests/*; do
            ID=$(basename $REQUEST)
            mysql -e "CREATE DATABASE $ID;"
            timeout 5 mysql $ID --table < $REQUEST 1>tmp/$ID 2>&1
            if [ $? -eq 124 ]; then
              echo -e "\n\n---------------------\nTimeout reached, aborting\n---------------------\n" >> tmp/$ID
            fi
            mysql -e "DROP DATABASE $ID;"
            mv tmp/$ID results/$ID
            rm $REQUEST
        done
    fi
    sleep 1
done
