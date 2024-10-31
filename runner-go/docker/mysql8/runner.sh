#!/bin/bash

cd /app

mkdir -p tmp

export MYSQL_PWD=mysql8root

while true; do
    mysql -u root -e "SHOW DATABASES;" 1>/dev/null 2>&1
    if [ $? -eq 0 ]; then
        echo "Starting mysql8 runner"
        break
    fi
    echo "Waiting for mysql8"
    sleep 2
done

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
            mysql -u root -e "CREATE DATABASE tmp_$ID;"
            timeout 5 mysql -u root tmp_$ID --table < $REQUEST 1>>tmp/$ID 2>&1
            if [ $? -eq 124 ]; then
              echo -e "\n\n-------------------------\nTimeout reached, aborting\n-------------------------\n" >> tmp/$ID
            fi
            mysql -u root -e "DROP DATABASE tmp_$ID;"
            rm $REQUEST
            mv tmp/$ID results/$ID
        done
    fi
    sleep 0.01
done
