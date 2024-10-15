#!/bin/bash

adduser --disabled-password restricted_user

mkdir -p go
chmod -R 755 go
mkdir -p tmp
chmod -R 700 tmp

while [ True ]; do
    if [ -n "$(ls go)" ]; then
      rm go/*
    fi
    if [ -n "$(ls tmp)" ]; then
      rm tmp/*
    fi
    if [ -n "$(ls requests)" ]; then
        for REQUEST in requests/*; do
            echo $REQUEST
            ID=$(basename $REQUEST)
            touch tmp/$ID
            chmod 700 tmp/$ID
            mv $REQUEST go/$ID.go
            chmod 755 go/$ID.go
            su -c "timeout 10 go run go/$ID.go" restricted_user 1>>tmp/$ID 2>&1
            if [ $? -eq 124 ]; then
              echo -e "\n\n-------------------------\nTimeout reached, aborting\n-------------------------\n" >> tmp/$ID
            fi
            rm go/$ID.go
            mv tmp/$ID results/$ID
        done
    fi
    sleep 0.01
done
