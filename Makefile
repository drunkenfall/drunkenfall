SOURCEDIR = .
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY = drunkenfall

VERSION = $(shell git describe --always --tags)
BUILDTIME = `date +%FT%T%z` # ISO-8601

LDFLAGS = -ldflags "-X main.version=${VERSION} -X main.buildtime=${BUILDTIME}"

export GOPATH := $(shell go env GOPATH)
.DEFAULT_GOAL: $(BINARY)

check: test lint


BINARY = drunkenfall
$(BINARY): $(SOURCES)
	go build -v ${LDFLAGS} -o ${BINARY}

.PHONY: $(BINARY)-start
$(BINARY)-start: $(BINARY)
	./$(BINARY)

.PHONY: dist
dist: $(BINARY)
	cd js; npm run build

.PHONY: install-linter
install-linter:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(GOPATH)/bin v1.12.2

.PHONY: test
test:
	GIN_MODE=test go test -v ./towerfall

.PHONY: cover
cover:
	go test -coverprofile=cover.out ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: npm
npm: js/package.json
	cd js; npm install

.PHONY: npm-start
npm-start:
	cd js; PORT=42002 npm run dev

.PHONY: npm-sass
npm-sass:
	cd js; npm rebuild node-sass

.PHONY: npm-dist
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
	sudo $(GOPATH)/bin/caddy

.PHONY: caddy-local
caddy-local: download-caddy
	sudo $(GOPATH)/bin/caddy -conf Caddyfile.local

.PHONY: postgres
postgres:
	docker-compose up postgres

.PHONY: psql
psql:
	@psql --host localhost --user postgres drunkenfall

DB := ./data/db.sql
DBPARAMS := --host localhost --user postgres

.PHONY: reset-db
reset-db:
	test -n "$(DRUNKENFALL_RESET_DB)"
	psql  -c "DROP DATABASE drunkenfall"
	psql $(DBPARAMS) -c "CREATE DATABASE drunkenfall"
	psql $(DBPARAMS) drunkenfall < $(DB)

.PHONY: reset-test-db
reset-test-db:
	psql $(DBPARAMS) -c "DROP DATABASE test_drunkenfall"
	psql $(DBPARAMS) -c "CREATE DATABASE test_drunkenfall"
	psql $(DBPARAMS) test_drunkenfall < $(DB)

.PHONY: make-test-db
make-test-db:
	pg_dump $(DBPARAMS) drunkenfall > $(DB)
