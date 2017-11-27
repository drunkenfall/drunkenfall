## -*- docker-image-name: "drunkenfall" -*-
FROM golang:latest

ENV DF_ROOT=/go/src/github.com/drunkenfall/drunkenfall/
WORKDIR $DF_ROOT

# I hate everything.
RUN curl -sL https://deb.nodesource.com/setup_6.x | bash -
RUN apt-get install nodejs

RUN go get github.com/kardianos/govendor

COPY js/package*.json ./js/
RUN cd js; npm install --only-production

COPY vendor/vendor.json ./vendor/
RUN govendor sync -v

# Copy mostly static js stuff
COPY js/build/ ./js/build/
COPY js/config/ ./js/config/
COPY js/.eslintrc.js ./js/.eslintrc.js
COPY js/.babelrc ./js/.babelrc
COPY js/static/ ./js/static/
COPY js/index.html ./js/

COPY Caddyfile .
COPY Makefile .

COPY websockets/ ./websockets/
COPY towerfall/ ./towerfall/
COPY js/src ./js/src/
COPY *.go ./
COPY .git ./.git

RUN make drunkenfall
RUN make npm-dist
