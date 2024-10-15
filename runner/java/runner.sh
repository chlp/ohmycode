#!/bin/bash

adduser --disabled-password --gecos "" restricted_user

mkdir -p java
chmod -R 755 java
mkdir -p tmp
chmod -R 700 tmp

while [ True ]; do
    if [ -n "$(ls java)" ]; then
      rm -rf java/*
    fi
    if [ -n "$(ls tmp)" ]; then
      rm -rf tmp/*
    fi
    if [ -n "$(ls requests)" ]; then
        for REQUEST in requests/*; do
            echo $REQUEST
            ID=$(basename $REQUEST)
            touch tmp/$ID
            chmod 700 tmp/$ID
            mkdir -p java/$ID
            mv $REQUEST java/$ID/Main.java
            chmod -R 755 java/$ID
            timeout 10 javac java/$ID/Main.java 1>tmp/$ID 2>&1
            if [ $? -eq 0 ]; then
                cd java/$ID
                su -c "timeout 10 java Main" restricted_user 1>>../../tmp/$ID 2>&1
                if [ $? -eq 124 ]; then
                    echo -e "\n\n-------------------------\nTimeout reached, aborting\n-------------------------\n" >> ../../tmp/$ID
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
