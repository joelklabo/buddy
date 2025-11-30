.PHONY: run build lint test fmt install

CONFIG ?= config.yaml
BIN ?= bin/nostr-codex-runner

run:
	CONFIG=$(CONFIG) go run ./cmd/runner

build:
	go build -o $(BIN) ./cmd/runner

lint:
	go vet ./...

test:
	go test ./...

fmt:
	gofmt -w cmd internal

install:
	go install ./cmd/runner

verify:
	@gofmt -l cmd internal | tee /tmp/gofmt.out
	@test ! -s /tmp/gofmt.out || (echo "gofmt needed" && cat /tmp/gofmt.out && false)
	go vet ./...
	@command -v staticcheck >/dev/null 2>&1 || go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./...
	go test -race -covermode=atomic ./...
