all: build

.PHONY: assets
assets:
		$$GOPATH/bin/go-bindata -o src/assets.go assets/...

.PHONY: build
build: assets
		go build
