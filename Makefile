.NOTPARALLEL: static-arg nowindowsgui build start
SHELL = sh

all: build
static: static-arg build
run: nowindowsgui build start

LDFLAGS = -s -w -X 'git.randomchars.net/FreeNitori/FreeNitori/nitori/state.version=$(shell echo -n `git describe --tags`; if ! [ "`git status -s`" = '' ]; then echo -n '-dirty'; fi)' -X 'git.randomchars.net/FreeNitori/FreeNitori/nitori/state.revision=$(shell git rev-parse --short HEAD)'

ifeq ($(shell go env GOOS), windows)
   WINDOW_LDFLAGS = -H windowsgui
   WINDOW_SYSO_REMOVE = rm -f cmd/server/freenitori.syso
   $(shell cp assets/freenitori.syso cmd/server/freenitori.syso)
endif

.PHONY: static-arg
static-arg:
	$(eval STATIC_LDFLAGS = -extldflags "-static")

.PHONY: image
image:
	CGO_ENABLED=0 go build -tags netgo -ldflags="-w -extldflags \"-static\"" -o build/init $$PWD/cmd/init
	sh assets/os/build.sh

.PHONY: nowindowsgui
nowindowsgui:
	$(eval WINDOW_LDFLAGS = )

.PHONY: build
build:
	@echo "Building FreeNitori..."
	@go build -tags=jsoniter -ldflags="$(LDFLAGS) $(STATIC_LDFLAGS) $(WINDOW_LDFLAGS)" -o build/freenitori$(shell go env GOEXE) $$PWD/cmd/server
	@echo "Building nitorictl..."
	@go build -tags=jsoniter -ldflags="$(LDFLAGS)" -o build/nitorictl$(shell go env GOEXE) $$PWD/cmd/cli
	@$(WINDOW_SYSO_REMOVE)

.PHONY: start
start:
	@./build/freenitori
