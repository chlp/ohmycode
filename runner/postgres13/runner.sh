#!/bin/bash

cd /app

export PGPASSWORD=password

while [ True ]; do
    if [ -n "$(ls tmp)" ]; then
      rm tmp/*
    fi
    if [ -n "$(ls requests)" ]; then
        for REQUEST in requests/*; do
            echo $REQUEST
            ID=$(basename $REQUEST)
            psql -U user -d mydatabase -c "CREATE DATABASE tmp_$ID;" 1>/dev/null 2>&1
            timeout 5 psql -U user -d tmp_$ID -q -c "\pset format wrapped" -f $REQUEST 1>tmp/$ID 2>&1
            if [ $? -eq 124 ]; then
              echo -e "\n\n-------------------------\nTimeout reached, aborting\n-------------------------\n" >> tmp/$ID
            fi
            psql -U user -d mydatabase -c "DROP DATABASE tmp_$ID;" 1>/dev/null 2>&1
            rm $REQUEST
            mv tmp/$ID results/$ID
            chmod 777 results/$ID
        done
    fi
    sleep 0.01
done
