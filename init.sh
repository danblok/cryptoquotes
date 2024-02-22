#!/bin/bash

# Creating a dir for keys
mkdir -p data

/bin/bash ${PWD}/scripts/genkeys.sh

docker-compose up --build -d
