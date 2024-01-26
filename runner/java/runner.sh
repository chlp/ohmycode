#!/bin/bash

mkdir -p java
mkdir -p tmp
while [ True ]; do
    if [ ! -z "$(ls java)" ]; then
      rm -rf java/*
    fi
    if [ ! -z "$(ls tmp)" ]; then
      rm -rf tmp/*
    fi
    if [ ! -z "$(ls requests)" ]; then
        for REQUEST in requests/*; do
            echo $REQUEST
            ID=$(basename $REQUEST)
            mkdir -p java/$ID
            mv $REQUEST java/$ID/Main.java
            javac java/$ID/Main.java 1>tmp/$ID 2>&1
            if [ $? -eq 0 ]; then
                cd java/$ID
                timeout 5 java Main 1>../../tmp/$ID 2>&1
                if [ $? -eq 124 ]; then
                    echo -e "\n\n-------------------------\nTimeout reached, aborting\n-------------------------\n" >> tmp/$ID
                fi
                cd ../..
                rm -rf java/$ID
            else
                echo -e "\n\n-------------------------\Compilation failed\n-------------------------\n" >> tmp/$ID
            fi
            mv tmp/$ID results/$ID
        done
    fi
    sleep 0.01
done
