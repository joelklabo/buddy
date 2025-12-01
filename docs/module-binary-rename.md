# Module & Binary Rename Plan – Issue 3oa.39

Current module: `module nostr-codex-runner`
Current binary: `nostr-codex-runner` (cmd/runner)
Target module: `module github.com/joelklabo/buddy`
Target binary: `buddy` (with alias symlink `nostr-buddy`; optional `bud` opt-in)

Steps
1) After repo rename (3oa.38), update `go.mod` module path and run `go mod tidy`.
2) Update import paths across repo (`github.com/joelklabo/nostr-codex-runner` → `github.com/joelklabo/buddy`).
3) Rename `cmd/runner` package output to `buddy`:
   - Adjust `Makefile` targets (`BIN` default).
   - Update `scripts/install.sh` artifact names.
   - Update any references in docs and sample commands.
4) Provide alias binary:
   - In goreleaser config, either build twice or add post-build symlink `nostr-buddy` (preferred) and optionally `bud` (off by default).
5) Backward compat:
   - Keep `nostr-codex-runner` binary name available via symlink for one release cycle; emit warning on startup pointing to `buddy`.
   - ENV var shims: keep reading `NCR_CONFIG` but log deprecation; new env name could be `BUDDY_CONFIG`.
6) Tests/CI:
   - Update workflow caches and binary names (artifacts in CI/release).
   - Ensure `go test ./...` passes with new module path; refresh coverage badges.
7) External references:
   - Update package docs badge, go reference badge, README commands, systemd template, recipes, sample configs.
8) Validation steps:
   - Fresh clone after rename; `go test ./...`.
   - `make build` produces `bin/buddy` and `bin/nostr-buddy` symlink.
   - `buddy version` shows new module path.

Risks
- Downstream imports break: mitigate with replace directive temporarily and clear release notes.
- Users with old binary names: provide shim and warning.
