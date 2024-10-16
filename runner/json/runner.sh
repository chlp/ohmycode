#!/bin/bash

adduser --disabled-password restricted_user

mkdir -p tmp
chmod -R 700 tmp

while [ True ]; do
    if [ -n "$(ls tmp)" ]; then
      rm tmp/*
    fi
    if [ -n "$(ls requests)" ]; then
        for REQUEST_FILEPATH in requests/*; do
            echo $REQUEST_FILEPATH
            ID=$(basename $REQUEST_FILEPATH)
            touch tmp/$ID
            chmod 700 tmp/$ID
            su -c "timeout 5 cat $REQUEST_FILEPATH | jq" restricted_user 1>>tmp/$ID 2>&1
            if [ $? -eq 124 ]; then
              echo -e "\n\n-------------------------\nTimeout reached, aborting\n-------------------------\n" >> tmp/$ID
            fi
            rm $REQUEST_FILEPATH
            mv tmp/$ID results/$ID
        done
    fi
    sleep 0.01
done
