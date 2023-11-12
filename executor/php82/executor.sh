#!/bin/bash

while [ True ]; do
    if [ ! -z "$(ls requests)" ]; then
        for REQUEST in requests/*; do
            ID=$(basename $REQUEST)
            php $REQUEST 1>results/$ID 2>&1
            rm $REQUEST
        done
    fi
    sleep 1
done
