#!/bin/bash

cd /app

export PGPASSWORD=password

while [ True ]; do
    if [ ! -z "$(ls tmp)" ]; then
      rm tmp/*
    fi
    if [ ! -z "$(ls requests)" ]; then
        for REQUEST in requests/*; do
            echo $REQUEST
            ID=$(basename $REQUEST)
            psql -U user -d mydatabase -c "CREATE DATABASE tmp_$ID;" 1>/dev/null 2>&1
            timeout 5 psql -U user -d tmp_$ID < $REQUEST 1>tmp/$ID 2>&1
            if [ $? -eq 124 ]; then
              echo -e "\n\n-------------------------\nTimeout reached, aborting\n-------------------------\n" >> tmp/$ID
            fi
            psql -U user -d mydatabase -c "DROP DATABASE tmp_$ID;" 1>/dev/null 2>&1
            mv tmp/$ID results/$ID
            rm $REQUEST
        done
    fi
    sleep 0.01
done
