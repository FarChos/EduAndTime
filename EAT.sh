#!/bin/zsh

docker-compose up -d

sleep 24

docker exec -it couchbase /init.sh 
sleep 30
docker exec -it couchbase /index.sh 
