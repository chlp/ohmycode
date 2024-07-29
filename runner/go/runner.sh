#!/bin/bash

mkdir -p go
mkdir -p tmp
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
            mv $REQUEST go/$ID.go
            timeout 10 go run go/$ID.go 1>tmp/$ID 2>&1
            if [ $? -eq 124 ]; then
              echo -e "\n\n-------------------------\nTimeout reached, aborting\n-------------------------\n" >> tmp/$ID
            fi
            rm go/$ID.go
            mv tmp/$ID results/$ID
            chmod 777 results/$ID
        done
    fi
    sleep 0.01
done
