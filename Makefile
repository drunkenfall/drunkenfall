SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=drunkenfall

VERSION=$(shell git describe --dirty --always --tags)
BUILDTIME=`date +%FT%T%z` # ISO-8601

LDFLAGS=-ldflags "-X drunkenfall.version=${VERSION} -X  drunkenfall.buildtime=${BUILDTIME}"

.DEFAULT_GOAL: all
.PHONY: download install test clean npm

all: clean download npm test $(BINARY)

$(BINARY): download $(SOURCES)
	go build -v ${LDFLAGS} -o ${BINARY}

download:
	go get -t -d -v ./...

install:
	go install -v ${LDFLAGS} ./...

test:
	go test -v

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

npm:
	cd js; npm install
