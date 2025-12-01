# Help/Man Examples – Issue 3oa.49

Goal: make `buddy help run` and `buddy help presets` show copy-paste examples.

Proposed snippets
- help run:
  - `buddy run mock-echo` (offline smoke) — expects echo replies in logs.
  - `buddy run claude-dm` — Nostr DM to Claude/OpenAI; needs API key + relays.
  - `buddy run config.yaml` — positional config path; honors BUDDY_CONFIG when flag omitted.
- help presets:
  - `buddy presets` — lists names + one-line summaries.
  - `buddy presets claude-dm` — shows summary + "Run: buddy run claude-dm" + hint to `--yaml` for full config.
  - `buddy presets claude-dm --yaml` — prints YAML.

Implementation plan
- Update `printHelp` to print examples for run/presets commands.
- Ensure man page includes these examples under COMMANDS/EXAMPLES.
- Tests: capture help output and assert example lines exist.
