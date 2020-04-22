#!/bin/bash
export GOOS=linux
export GOARCH=amd64
export ESM_LISTEN=${1:-8080}
echo "\n $ESM_LISTEN"
cd ~/external-service-marketplace
rm esmd
make esmd
docker build . -t esmd:test --build-arg PORT=${ESM_LISTEN}
echo "\n MADE THE IMAGE"
docker run -p ${ESM_LISTEN}:${ESM_LISTEN} esmd:test 
