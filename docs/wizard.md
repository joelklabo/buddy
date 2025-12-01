# buddy wizard – Issue 3oa.22

Quickest way to produce a working config without editing YAML.

## Usage
```bash
nostr-codex-runner wizard [config-path]
```
- If no path is given, writes `~/.config/buddy/config.yaml` (creates parent dirs).
- If the file exists, prompts before overwrite.

## What it asks
1) Relays (defaults to two public relays)
2) Nostr private key (masked input)
3) Allowed pubkeys (comma-separated)
4) Agent choice: http, copilotcli, or echo
5) Whether to enable shell action (warned as high risk)
6) Whether to dry-run (preview without writing)

## What it writes
- Transport: nostr with your relays and keys
- Agent: chosen type
- Actions: readfile always; shell if enabled
- Projects: default project at `.`
- Storage: `~/.local/share/buddy/state.db`

## After running
- Output shows the config path.
- Start the runner:
  ```bash
  nostr-codex-runner run -config ~/.config/buddy/config.yaml
  ```

## Notes
- Uses masked prompts; secrets are not echoed.
- Respects legacy config/search paths via CLI defaults, but wizard writes to the buddy path by default.
- Add `--config` to `run` if you wrote to a custom path.
- Future: will surface presets and mock transport option.
