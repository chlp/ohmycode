#!/bin/bash

adduser --disabled-password restricted_user

mkdir -p tmp
while [ True ]; do
    if [ -n "$(ls tmp)" ]; then
      rm tmp/*
    fi
    if [ -n "$(ls requests)" ]; then
        for REQUEST in requests/*; do
            echo $REQUEST
            ID=$(basename $REQUEST)
            su -c "timeout 5 php $REQUEST 1>tmp/$ID 2>&1" restricted_user
            if [ $? -eq 124 ]; then
              echo -e "\n\n-------------------------\nTimeout reached, aborting\n-------------------------\n" >> tmp/$ID
            fi
            rm $REQUEST
            mv tmp/$ID results/$ID
            chmod 777 results/$ID
        done
    fi
    sleep 0.01
done
