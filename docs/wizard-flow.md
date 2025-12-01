# Wizard Question Flow – Issue 3oa.18

This is the draft script/IA for `buddy wizard [config-path]`.

## Steps
1) **Welcome**
   - Explain that it will write a config (default: `~/.config/buddy/config.yaml`).
   - Confirm overwrite if file exists.

2) **Choose a preset or blank**
   - Options: `claude-dm`, `nostr-copilot-shell`, `local-llm`, `mock-echo`, `blank`.
   - Selecting a preset pre-fills later answers; user can still edit inputs.

3) **Transport**
   - Default: `nostr`.
   - If nostr: ask relays (prefill e.g., `wss://relay.damus.io,wss://nos.lol`), ask for private key (masked), and allowed pubkeys (comma list). Validate hex length.
   - Offer `mock` transport for offline; skip secrets in that path.

4) **Agent** (may be set by preset)
   - Options: `http-claude/openai`, `copilotcli`, `local-http` (specify URL), `echo` (for mock). Collect:
     - For HTTP: provider name, API key (masked), base URL (optional), model name, timeout.
     - For copilotcli: binary path (default `copilot`), working dir, extra args? keep minimal.
     - For local: endpoint URL or binary command.

5) **Actions**
   - Default none. Offer: `shell` (warn, require confirmation), `readfile`, `writefile`.
   - For shell: ask workdir (default `.`), timeout, max_output.
   - For read/write: roots list.

6) **Runner settings**
   - Session timeout (default 60m), initial prompt (optional), log level (info|debug), storage path (BoltDB default `~/.local/share/buddy/state.db`).

7) **Summary + path**
   - Show derived preset name if any, transport/agent/actions summary, and config path.
   - Ask "Write this config now?" (yes/no). If no, abort without writing.

8) **Write config**
   - Create parent dirs; write YAML.

9) **Next steps prompt**
   - Offer to run now: `buddy run <preset>` or `buddy run <config-path>`.
   - If yes, start runner and stream logs until Ctrl+C.

## Copy snippets (examples)
- Welcome: "We'll ask a few questions and write a YAML config. Secrets stay local and are masked as you type."
- Shell warning: "Shell actions let the agent run commands on your machine. Enable only for trusted operators."
- Collision note: "If you already have another `buddy` on PATH (Buddy.Works), consider using the `nostr-buddy` alias."

## Validation/UX rules
- Mask secret inputs (private key, API keys).
- Validate hex keys lengths for nostr; retry loop with clear error.
- For overwrites, require explicit "yes".
- Keep prompts single-line; avoid flags; no colors required for MVP.

## Outputs
- Config YAML written to chosen path.
- Summary printed: path, transport, agent, actions, preset used/overridden.
- Optional: write a `.wizard.yaml` scratch for debugging (off by default).

## Open questions
- Do we auto-add default relays, or force user to enter at least one?
- Should we collect metrics endpoint and health port? (Probably skip in MVP to stay short.)
- Allow `--dry-run` and `--no-run`? lean to `--dry-run` only if needed for tests.
