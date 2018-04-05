#
# go-fibo 
# Author: Damian Janiszewski
#
# Multistage Dockerfile definition
#
# Stage 0: golang compiler and build container
FROM golang:latest as builder-go
LABEL author=damian-janiszewski

ENV GOLANG_VERSION 1.9.2
WORKDIR /go/src/

# Get dependencies
RUN go get -v -d -tags static github.com/gorilla/mux

# Copy sources
COPY *.go /go/src/
ENV CGO_ENABLED=0
RUN go build -v go-fibo.go

# Stage 1: running container 
FROM alpine:latest

LABEL version="0.02"
LABEL author=damian-janiszewski

# Copy binaries from stage 0 builder container
COPY --from=builder-go /go/src/go-fibo /usr/local/bin/

RUN mkdir -p /usr/local/bin/ \
        && chmod ug+x /usr/local/bin/go-fibo

CMD ["/usr/local/bin/go-fibo"]
