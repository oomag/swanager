BINARY=swanager
VERSION=1.0.0
BUILD_TIME=`date +%FT%T%z`
LDFLAGS=-ldflags "-linkmode external -extldflags -static -X github.com/da4nik/swanager/core.Version=${VERSION} -X github.com/da4nik/swanager/core.BuildTime=${BUILD_TIME}"

SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

.PHONY: build run install clean
.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
	glide install
	go build ${LDFLAGS} -o ${BINARY} swanager.go

build: $(BINARY)

run:
	go run ${BINARY}.go

install:
	go install ${LDFLAGS} ./...

clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
