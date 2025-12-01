# Local Dev Setup (Contributors)

Scope: for contributors. End users should install the `buddy` binary from releases or brew.

Prereqs

- Go ≥1.24
- Copilot CLI or HTTP keys optional (only if you test those agents)
- `bd` installed (issue tracking)
- Access to test Nostr relays (for integration flows)

Setup steps

1. Clone:

   ```bash
   git clone https://github.com/joelklabo/buddy.git
   cd buddy
   ```

2. Deps:

   ```bash
   go mod download
   ```

3. Build & run:

   ```bash
   make build      # bin/buddy
   make run        # uses config.yaml
   ```

4. Formatting & lint:

   ```bash
   make fmt
   make lint       # go vet ./...
   ```

5. Tests:

   ```bash
   go test ./...
   ```

6. Sample configs:

   - `config.example.yaml` → copy to `config.yaml`
   - `sample-flows/` has preset configs (copilot-shell, etc.).

Paths & state

- Config: `./config.yaml` for repo work; default runtime path `~/.config/buddy/config.yaml`.
- State DB: `state.db` (BoltDB) — excluded by .gitignore; keep test DBs out of git.
- Presets: embedded under `internal/presets/data`; user overrides in `~/.config/buddy/presets/`.

Testing notes

- Prefer table-driven tests; place alongside code in `_test.go`.
- For nostr-dependent tests, use mock transport or inject relay URLs via test config.

bd workflow

- One issue per commit; create under epic `buddy-3oa`.
- Commit message should close the issue, e.g., `Some change (closes buddy-3oa.X)`.

Module

- Path: `github.com/joelklabo/buddy`.
- Binary: `buddy` (no shipped alias).
