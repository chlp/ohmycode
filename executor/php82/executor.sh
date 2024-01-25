#!/bin/bash

mkdir tmp
while [ True ]; do
    if [ ! -z "$(ls tmp)" ]; then
      rm tmp/*
    fi
    if [ ! -z "$(ls requests)" ]; then
        for REQUEST in requests/*; do
            ID=$(basename $REQUEST)
            timeout 5 php $REQUEST 1>tmp/$ID 2>&1
            if [ $? -eq 124 ]; then
              echo -e "\n\n-------------------------\nTimeout reached, aborting\n-------------------------\n" >> tmp/$ID
            fi
            mv tmp/$ID results/$ID
            rm $REQUEST
        done
    fi
    sleep 0.01
done
