# Security & Secrets – Issue 3oa.14

Principles
- Least privilege: enable only needed transports/actions; keep allowlists tight.
- No secret logging: never print private keys, API tokens, or DM content in logs.
- Operators only: treat shell action as high risk; keep access to trusted pubkeys.

Secrets handling
- Nostr keys and API tokens should be provided via config/wizard prompts; not committed to git.
- Wizard must mask secret input and avoid writing secrets to stdout/stderr.
- If env vars are later supported, document them as opt-in and warn about process listing.

Logging
- Default logs at info level are non-verbose. For debug, advise against running on shared hosts if secrets may surface.
- Prefer JSON logs only when needed; redact where possible.

Transport trust
- Use trusted relays; prefer private relays for production. Avoid sending secrets in plaintext messages.
- Configure `runner.allowed_pubkeys` to restrict control; keep list short and reviewed.

Actions safety
- `shell`: require explicit enable; set `timeout_seconds` and `max_output`; avoid broad workdirs.
- `readfile`/`writefile`: set `roots` allowlist; avoid `/` or user home unless necessary.

Sandbox/approvals
- Codex sandbox flags (`sandbox: danger-full-access`, `approval: never`) should be used only on trusted machines; document risks in README/FAQ.

Storage
- BoltDB state file may contain session/thread IDs and prompts. Place it in trusted locations; avoid world-readable permissions.

Release artifacts
- Ship checksums; consider signing (cosign) for verification.

Incident response
- Rotate keys if leaked; update config and restart runner.
- Remove compromised relays from config; review logs for abuse.
