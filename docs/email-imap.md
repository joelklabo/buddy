# IMAP/SMTP Transport (beta)

Status: beta/skeleton. Use Mailgun webhook for production; IMAP/SMTP still needs live testing.

Config example:

```yaml
transports:
  - type: email
    config:
      mode: imap
      id: email-imap
      host: imap.example.com
      port: 993
      username: inbox@example.com
      password: $IMAP_PASSWORD
      folder: INBOX
      idle: true
      smtp_host: smtp.example.com
      smtp_port: 587
      allow_senders: ["alice@example.com"]
      max_bytes: 262144
```

Notes/limitations:

- Plain IMAP/SMTP; no OAuth2 yet.
- Processes text bodies only; attachments/HTML ignored.
- Allowlist and size caps enforced.
- Health endpoint not wired; add if you deploy this path.

Testing plan (manual):

- Use a test mailbox; send an email to inbox, watch logs for inbound, ensure replies delivered with threading.
- Add race-free tests with a mock IMAP server (TODO).

For production, prefer Mailgun inbound webhooks until this path is hardened and covered by tests.
