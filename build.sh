#!/bin/bash

docker run --rm -v "$PWD":/usr/local/go/src/github.com/da4nik/swanager -w /usr/local/go/src/github.com/da4nik/swanager golang:1.7.5 bash -c 'go get -v -d && go build -o ./swanager  -v -ldflags "-linkmode external -extldflags -static"'

docker build -t swanager .

rm -f ./swanager
