# Contributing

Thanks for helping improve the project! This repo uses `bd` for issue tracking and expects one commit per closed issue.

## How we work
- Create or pick an issue under epic `nostr-codex-runner-2zo` using `bd create --parent nostr-codex-runner-2zo ...`.
- Keep work focused: one logical change per issue; close the issue with a commit that references it.
- Prefer small PRs with context and testing notes.

## Development setup
- Requirements: Go ≥1.22, Codex CLI on PATH, `bd` installed, access to test Nostr relays.
- Install deps: `go mod download`
- Format: `gofmt -w ./cmd ./internal`
- Lint: `go vet ./...`
- Tests: `go test ./...`
- Run: `make run` (expects `config.yaml`)

## Commit and PR guidelines
- Reference the `bd` issue in the commit subject, e.g. `Add metrics exporter (closes nostr-codex-runner-2zo.3)`.
- Keep commit messages in imperative mood; avoid squashing unrelated work together.
- Include a short testing section in PR descriptions (commands run, notable results).

## Coding conventions
- Go style: idiomatic, small functions, explicit error wrapping, prefer `context.Context` plumbed.
- Config and secrets: never commit real keys; `.gitignore` already excludes `config.yaml` and state files.
- Logging: use structured logs (`slog`) and avoid leaking secret material.

## Security disclosures
See `SECURITY.md` for how to report vulnerabilities. Please **do not** open public issues for security findings.
