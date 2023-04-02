repo=$${REPO:-damianjaniszewski/go-fibo}
version=$${VERSION:-0.0.16}
tag=$(version)

build:
	go build -v go-fibo.go config.go handlers.go
	# CGO_ENABLED=0 go build -v -a -tags 'static netgo' -ldflags '-w' ./go-fibo.go ./config.go ./handlers.go

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -buildmode=exe -tags 'static netgo' -ldflags '-w' go-fibo.go config.go handlers.go
		
run:
	go build -v go-fibo.go config.go handlers.go
	# CGO_ENABLED=0 go build -v -a -tags 'static netgo' -ldflags '-w' ./go-fibo.go ./config.go ./handlers.go
	./go-fibo

build-container:
	tar -czv -f context.tar.gz ./Dockerfile ./go.mod ./go.sum ./*.go
	docker rmi $(repo):$(tag) $(repo):latest || true
	docker build -t $(repo):$(tag) -t $(repo):latest .
	docker push $(repo):$(tag)
	docker push $(repo):latest
	rm context.tar.gz || true

ver:
	@echo $(repo):$(tag)
