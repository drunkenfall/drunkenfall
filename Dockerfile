FROM golang:alpine as builder

ENV DF_ROOT=/go/src/gitlab.com/one-eye/drunkenfall/
WORKDIR $DF_ROOT

RUN apk --no-cache add curl alpine-sdk
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

COPY Gopkg.lock Gopkg.toml ./
RUN dep ensure -v -vendor-only

COPY Makefile ./

COPY faking/ ./faking/
COPY towerfall/ ./towerfall/
COPY *.go ./
# For the versioning to work
COPY .git ./.git

RUN make install

FROM scratch
COPY --from=builder /go/bin/drunkenfall /
ENTRYPOINT ["/drunkenfall"]
