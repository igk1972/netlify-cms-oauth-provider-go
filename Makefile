TARGET_BIN = netlify-cms-oauth-provider
TARGET_ARCH = amd64
SOURCE_MAIN = main.go
LDFLAGS = -s -w

all: build

build: build-darwin build-linux build-windows

build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=$(TARGET_ARCH) go build -ldflags "$(LDFLAGS)" -o bin/$(TARGET_BIN)_darwin-amd64 $(SOURCE_MAIN)

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=$(TARGET_ARCH) go build -ldflags "$(LDFLAGS)" -o bin/$(TARGET_BIN)_linux-amd64 $(SOURCE_MAIN)

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=$(TARGET_ARCH) go build -ldflags "$(LDFLAGS)" -o bin/$(TARGET_BIN)_windows-amd64.exe $(SOURCE_MAIN)

build-linux-noarch:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "$(LDFLAGS)" -o bin/$(TARGET_BIN) $(SOURCE_MAIN)

docker-build:
	docker build -t rubemlrm/netlify-cms-auth-provider:latest .

start:
	go run $(SOURCE_MAIN)
