package presets

import "fmt"

// List returns preset names and descriptions.
func List() map[string]string {
	return map[string]string{
		"claude-dm":           "Nostr DM to Claude/OpenAI HTTP agent (no shell by default)",
		"nostr-copilot-shell": "Nostr DM to Copilot CLI with shell action (trusted)",
		"local-llm":           "Nostr DM to local HTTP LLM endpoint",
		"mock-echo":           "Offline mock transport + echo agent",
	}
}

// Get returns the raw YAML for a preset, or an error if unknown.
func Get(name string) ([]byte, error) {
	switch name {
	case "claude-dm":
		return ClaudeDM, nil
	case "nostr-copilot-shell":
		return NostrCopilotShell, nil
	case "local-llm":
		return LocalLLM, nil
	case "mock-echo":
		return MockEcho, nil
	default:
		return nil, fmt.Errorf("unknown preset %s", name)
	}
}
