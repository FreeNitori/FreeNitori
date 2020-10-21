all: build
test: build-test run-test

.PHONY: deps
deps:
	go get -u github.com/go-bindata/go-bindata/...

.PHONY: assets
assets:
	$$GOPATH/bin/go-bindata -o nitori/config/assets.go -pkg config -prefix assets/ ./assets/* ./assets/web/templates/*
	$$GOPATH/bin/go-bindata -fs -o nitori/web/static.go -pkg web -prefix assets/web/static/ ./assets/web/static/...

.PHONY: run-test
run-test:
	./build/freenitori

.PHONY: build
build: deps assets
	go build -tags=jsoniter -ldflags="-s -w" -o build/freenitori-supervisor proc/supervisor/main.go
	go build -tags=jsoniter -ldflags="-s -w" -o build/freenitori-chatbackend proc/chatbackend/main.go
	go build -tags=jsoniter -ldflags="-s -w" -o build/freenitori-webserver proc/webserver/main.go
	go build -tags=jsoniter -ldflags="-s -w" -o build/freenitori-console proc/console/main.go
	cp build/freenitori-supervisor build/freenitori

.PHONY: build-test
build-test: assets
	go build -tags=jsoniter -o build/freenitori-supervisor proc/supervisor/main.go
	go build -tags=jsoniter -o build/freenitori-chatbackend proc/chatbackend/main.go
	go build -tags=jsoniter -o build/freenitori-webserver proc/webserver/main.go
	go build -tags=jsoniter -o build/freenitori-console proc/console/main.go
	cp build/freenitori-supervisor build/freenitori
