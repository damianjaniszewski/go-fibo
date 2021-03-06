#
# go-fibo 
# Author: Damian Janiszewski
#
# Multistage Dockerfile definition
#
# Stage 0: golang compiler and build container
FROM golang:latest as builder-go

# ENV GOLANG_VERSION 1.9.2
WORKDIR /go/src/

# Copy sources
COPY *.go /go/src/

# Get dependencies
RUN go get -v -d -tags 'static netgo' "github.com/gofrs/uuid" \
	"github.com/gorilla/mux" \
	"github.com/prometheus/client_golang/prometheus" \
	"github.com/prometheus/client_golang/prometheus/promhttp" \
	"go.opencensus.io/zpages" \
	"github.com/sirupsen/logrus" \
	"github.com/codingconcepts/env" \
	"github.com/joho/godotenv" \
	"github.com/damianjaniszewski/zpages" && CGO_ENABLED=0 GOOS=linux go build -v -a -tags 'static netgo' -ldflags '-w' go-fibo.go config.go fibo.go

# Stage 1: running container 
FROM scratch

LABEL version="0.0.14"
LABEL author "Damian Janiszewski"

# Copy binaries from stage 0 builder container
COPY --from=builder-go /go/src/go-fibo /usr/local/bin/

#RUN mkdir -p /usr/local/bin/ \
#        && chmod ug+x /usr/local/bin/go-fibo

CMD ["/usr/local/bin/go-fibo"]
