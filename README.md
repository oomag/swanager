# swanager

## Build

Simple:
```bash
./build.sh
```

Full:

```bash
TAG=latest DELETE_BUILD_IMAGE=0 GOLANG_BUILD_IMAGE=golang:1.8 ./build.sh
```

All params are optional.

Param | Default value | Description
-----|-----|-----
TAG | latest | Tag of resulting image `swanager:TAG`
DELETE_BUILD_IMAGE | 1 | Whether of not, delete golang build image
GOLANG_BUILD_IMAGE | golang:1.8 | Golang build image


## Run

### Docker container

```bash
docker run -d -v /var/run/docker.sock:/var/run/docker.sock swanager
```

Required mounted resource is a docker socket to manage docker. 

Configure swanager container:

Env vars | Default value | Description
-----|-----|-----
SWANAGER_PORT | 4945 | API port
SWANAGER_LOG | stdout | Logfile
SWANAGER_MONGO_URL | mongodb://127.0.0.1:27017/swanager | mongodb url
SWANAGER_DB_NAME | swanager | Mongodb database name
SWANAGER_PATH_PREFIX | /data | Service mount points root
SWANAGER_LOCAL_SECRET_KEY | - | Secret key, to authenticate local services (if none, won't be authenticated) 


### Docker stack file

You may want to store mongodb to mounted volume, just check `swanager.yml` 

Start
```bash
docker stack deploy --compose-file swanager.yml swanager
```

Stop
```bash
docker stack rm swanager
```
