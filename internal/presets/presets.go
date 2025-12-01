package presets

import _ "embed"

//go:embed data/claude-dm.yaml
var ClaudeDM []byte

//go:embed data/copilot-shell.yaml
var CopilotShell []byte

//go:embed data/local-llm.yaml
var LocalLLM []byte

//go:embed data/mock-echo.yaml
var MockEcho []byte
