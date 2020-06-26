all: build

.PHONY: assets
assets:
		$$GOPATH/bin/go-bindata -o nitori/utils/assets.go -pkg utils assets/...

.PHONY: build
build: assets
		go build
