package presets

import (
	_ "embed"

	"github.com/joelklabo/buddy/internal/config"
)

//go:embed data/claude-dm.yaml
var ClaudeDM []byte

//go:embed data/copilot-shell.yaml
var CopilotShell []byte

//go:embed data/local-llm.yaml
var LocalLLM []byte

//go:embed data/mock-echo.yaml
var MockEcho []byte

// PresetDeps returns declared prerequisites for built-in presets.
func PresetDeps() map[string][]config.Dep {
	return map[string][]config.Dep{
		"copilot-shell": {
			{Name: "copilot", Type: "binary", Hint: "Install GitHub Copilot CLI: https://github.com/github/copilot-cli"},
		},
		"claude-dm": {
			{Name: "curl", Type: "binary", Optional: true, Hint: "Used for simple HTTP checks"},
		},
		"local-llm": {
			{Name: "curl", Type: "binary", Optional: true, Hint: "Useful for hitting local endpoints"},
		},
	}
}
