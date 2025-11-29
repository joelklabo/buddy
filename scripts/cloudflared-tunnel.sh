#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."

if ! command -v cloudflared >/dev/null 2>&1; then
  echo "cloudflared not found. Install via 'brew install cloudflared' or download from https://github.com/cloudflare/cloudflared/releases" >&2
  exit 1
fi

UI_ADDR=${UI_ADDR:-127.0.0.1:8080}

# Start a quick tunnel to the UI address
cloudflared tunnel --url "http://${UI_ADDR}"
