SOURCE_MAIN = main.go

all: build

build: goreleaser build --clean --snapshot

release: goreleaser release --clean

start:
	go run $(SOURCE_MAIN)
