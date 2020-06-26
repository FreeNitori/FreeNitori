all: build

.PHONY: assets
assets:
		$$GOPATH/bin/go-bindata -o nitori/config/assets.go -pkg config assets/...

.PHONY: build
build: assets
		go build
