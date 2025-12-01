package presets

import _ "embed"

//go:embed data/claude-dm.yaml
var ClaudeDM []byte

//go:embed data/nostr-copilot-shell.yaml
var NostrCopilotShell []byte

//go:embed data/local-llm.yaml
var LocalLLM []byte

//go:embed data/mock-echo.yaml
var MockEcho []byte
