# Logging Safety Guidelines

- Do not log secrets (private keys, API tokens, OAuth codes, emails bodies). Transport handlers must avoid logging raw payloads; log metadata only.
- Redact or omit: `private_key`, `api_key`, `token`, `authorization`, email `body`.
- Use structured logs with clear fields; avoid dumping entire request structs.
- Keep health endpoints free of sensitive data (status/version only).
- In tests, avoid printing sample secrets to stdout/stderr; prefer fixtures.

Checklist for new code:

- [ ] No secrets in slog fields or errors.
- [ ] Truncate large payloads.
- [ ] Add unit tests if a redaction helper is introduced.
