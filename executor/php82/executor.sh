#!/bin/bash

mkdir tmp
while [ True ]; do
    if [ ! -z "$(ls tmp)" ]; then
      rm tmp/*
    fi
    if [ ! -z "$(ls requests)" ]; then
        for REQUEST in requests/*; do
            ID=$(basename $REQUEST)
            php $REQUEST 1>tmp/$ID 2>&1
            mv tmp/$ID results/$ID
        done
    fi
    sleep 1
done
