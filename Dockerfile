#
# go-fibo 
# Author: Damian Janiszewski
#
# Multistage Dockerfile definition
#
# Stage 0: golang compiler and build container
FROM golang:latest as builder-go

WORKDIR /go/src/

# Copy sources and compile
RUN env
COPY go.mod go.sum /go/src/
RUN go mod download -json

COPY *.go /go/src/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -buildmode=exe -tags 'static netgo' -ldflags '-w' go-fibo.go config.go handlers.go

# Stage 1: running container 

# Stage 1: running container 
FROM scratch

LABEL version="0.0.16"
LABEL author "Damian Janiszewski"

# Copy binaries from stage 0 builder container
COPY --from=builder-go /go/src/go-fibo /usr/local/bin/

CMD ["/usr/local/bin/go-fibo"]
