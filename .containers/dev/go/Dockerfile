FROM golang:1.21.1-alpine3.17

ENV GOPATH /go
ENV GO111MODULE on

RUN apk update \
    && apk --no-cache add git make \
    && go install github.com/cosmtrek/air@v1.45.0
