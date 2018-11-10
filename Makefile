SOURCEDIR = .
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY = drunkenfall

VERSION = $(shell git describe --always --tags)
BUILDTIME = `date +%FT%T%z` # ISO-8601

LDFLAGS = -ldflags "-X main.version=${VERSION} -X main.buildtime=${BUILDTIME}"

export GOPATH := $(shell go env GOPATH)
# export PATH := $(GOPATH)/bin:$(PATH)
# gotype is disabled since it seems pointless and also produces 250 errors
# about not finding dependencies that definitely exists.
LINTER_ARGS = -j 4 \
  --enable-gc \
  --enable=gofmt \
  --enable=misspell \
  --enable=unparam \
  --enable=unused \
  --disable=errcheck \
  --disable=gotype \
  --deadline=10m \
  --tests

.DEFAULT_GOAL: all

.PHONY: clean clobber download install install-linter test cover race lint npm npm-dist caddy

all: clean download npm test race lint $(BINARY)

check: test lint

clean:
	rm -f $(BINARY)

clobber: clean
	rm -rf js/node_modules

BINARY = drunkenfall
$(BINARY): $(SOURCES)
	go build -v ${LDFLAGS} -o ${BINARY}

.PHONY: $(BINARY)-start
$(BINARY)-start: $(BINARY)
	./$(BINARY)

.PHONY: dist
dist: $(BINARY)
	cd js; npm run build

download:
	go get -t -d -v ./...

install:
	go install -v ${LDFLAGS} ./...

install-linter:
	go get -v -u github.com/alecthomas/gometalinter
	gometalinter --install

test:
	GIN_MODE=test go test -v ./towerfall -failfast

cover:
	go test -coverprofile=cover.out ./...

race:
	go test -race -v ./...

lint: install-linter
	gometalinter $(LINTER_ARGS) $(SOURCEDIR)

npm: js/package.json
	cd js; npm install

npm-start: npm
	cd js; PORT=42002 npm run dev

npm-dist: npm
	cd js; npm run build

.PHONY: vendor
vendor:
	dep ensure -v

.PHONY: docker
docker:
	docker-compose build

.PHONY: download-caddy
download-caddy:
	go get github.com/mholt/caddy/caddy
	go get github.com/caddyserver/builds
	cd $(GOPATH)/src/github.com/mholt/caddy/caddy; go run build.go

.PHONY: caddy
caddy: download-caddy
	sudo /home/thiderman/bin/caddy

.PHONY: postgres
postgres:
	docker-compose up postgres

.PHONY: psql
psql:
	@psql --host localhost --user postgres drunkenfall

.PHONY: reset-db
reset-db:
	[[ -n "$(DRUNKENFALL_RESET_DB)" ]] && \
	psql --host localhost --user postgres drunkenfall < init-db.sql \
    || echo "need to set DRUNKENFALL_RESET_DB"

.PHONY: reset-test-db
reset-test-db:
	psql --host localhost --user postgres -c "DROP DATABASE test_drunkenfall"
	psql --host localhost --user postgres -c "CREATE DATABASE test_drunkenfall"
	psql --host localhost --user postgres test_drunkenfall < test-db.sql

.PHONY: make-test-db
make-test-db:
	pg_dump --user postgres --host localhost drunkenfall > test-db.sql
