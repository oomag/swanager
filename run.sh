#!/bin/bash

docker run -d --name swanager -v /var/run/docker.sock:/var/run/docker.sock -e SWANAGER_MONGO_URL=mongodb://172.17.0.2/swanager -p 4945:4945  swanager
