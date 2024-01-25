#!/bin/bash

mkdir go
mkdir tmp
while [ True ]; do
    if [ ! -z "$(ls go)" ]; then
      rm go/*
    fi
    if [ ! -z "$(ls tmp)" ]; then
      rm tmp/*
    fi
    if [ ! -z "$(ls requests)" ]; then
        for REQUEST in requests/*; do
            ID=$(basename $REQUEST)
            mv $REQUEST go/$ID.go
            timeout 5 go run go/$ID.go 1>tmp/$ID 2>&1
            if [ $? -eq 124 ]; then
              echo -e "\n\n-------------------------\nTimeout reached, aborting\n-------------------------\n" >> tmp/$ID
            fi
            mv tmp/$ID results/$ID
            rm go/$ID.go
        done
    fi
    sleep 1
done
