all: build
test: build-test run

.PHONY: deps
deps:
	go get -u github.com/go-bindata/go-bindata/...

.PHONY: assets
assets:
	$$GOPATH/bin/go-bindata -o nitori/config/assets.go -pkg config -prefix assets/ ./assets/* ./assets/web/templates/*
	$$GOPATH/bin/go-bindata -fs -o nitori/web/static.go -pkg web -prefix assets/web/static/ ./assets/web/static/...

.PHONY: run
run:
	./FreeNitori

.PHONY: build
build: deps assets
	go build -ldflags="-s -w"
	upx --brute FreeNitori

.PHONY: build-test
build-test: assets
	go build
