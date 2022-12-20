#!/bin/bash

docker image build -t bogd/go-webserver .
docker container run -d bogd/go-webserver

sleep 1

echo "Testing Hello handler..."
curl 172.17.0.2:8080/hello
echo

echo "Testing static webpage..."
curl 172.17.0.2:8080/
echo

docker rm -f $(docker ps -aq)

