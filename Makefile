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
	./FreeNitoriTest

.PHONY: build
build: deps assets
	go build -tags=jsoniter -ldflags="-s -w"
	upx --best --color --brute FreeNitori

.PHONY: build-test
build-test: assets
	go build -tags=jsoniter -o FreeNitoriTest
