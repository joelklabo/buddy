.PHONY: run build lint test fmt install tunnel

CONFIG ?= config.yaml
BIN ?= bin/nostr-codex-runner
UI_ADDR ?= 127.0.0.1:8080

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

tunnel:
	UI_ADDR=$(UI_ADDR) ./scripts/cloudflared-tunnel.sh
