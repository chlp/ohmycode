#!/bin/bash

cd /app

mkdir -p tmp

export PGPASSWORD=password

while [ True ]; do
    if [ -n "$(ls tmp)" ]; then
      rm tmp/*
    fi
    if [ -n "$(ls requests)" ]; then
        for REQUEST in requests/*; do
            echo $REQUEST
            ID=$(basename $REQUEST)
            touch tmp/$ID
            chmod 744 tmp/$ID
            psql -U user -d mydatabase -c "CREATE DATABASE tmp_$ID;" 1>/dev/null 2>&1
            psql -U user -d tmp_$ID -c "CREATE USER tmp_user_$ID WITH PASSWORD 'password';" 1>/dev/null 2>&1
            psql -U user -d tmp_$ID -c "GRANT ALL PRIVILEGES ON DATABASE tmp_$ID TO tmp_user_$ID;" 1>/dev/null 2>&1
            timeout 5 psql -U tmp_user_$ID -d tmp_$ID -q -c "\pset format wrapped" -f $REQUEST 1>>tmp/$ID 2>&1
            if [ $? -eq 124 ]; then
              echo -e "\n\n-------------------------\nTimeout reached, aborting\n-------------------------\n" >> tmp/$ID
            fi
            psql -U user -d mydatabase -c "DROP DATABASE tmp_$ID;" 1>/dev/null 2>&1
            rm $REQUEST
            mv tmp/$ID results/$ID
        done
    fi
    sleep 0.01
done
