#!/bin/bash

TAG=${TAG:=latest}
GOLANG_BUILD_IMAGE=golang:1.8
DELETE_BUILD_IMAGE=${DELETE_BUILD_IMAGE:=1}


docker run --rm -v "$PWD":/usr/local/go/src/github.com/da4nik/swanager \
           -w /usr/local/go/src/github.com/da4nik/swanager \
           $GOLANG_BUILD_IMAGE \
           bash -c 'go get -v -d && go build -o ./swanager  -v -ldflags "-linkmode external -extldflags -static"'

if [ "$DELETE_BUILD_IMAGE" = 1 ]
then
    docker rmi $GOLANG_BUILD_IMAGE
fi

docker build -t swanager:$TAG .

rm -f ./swanager
