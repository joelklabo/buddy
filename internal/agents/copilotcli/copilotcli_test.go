package copilotcli

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"nostr-codex-runner/internal/core"
)

// create a fake gh binary that returns static text.
func writeFakeGh(t *testing.T, dir string) string {
	t.Helper()
	path := filepath.Join(dir, "gh")
	script := "#!/usr/bin/env bash\necho \"hi from copilot\"\n"
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatalf("write fake gh: %v", err)
	}
	return path
}

func TestCopilotAgentGenerates(t *testing.T) {
	td := t.TempDir()
	gh := writeFakeGh(t, td)

	ag := New(Config{Binary: gh})
	resp, err := ag.Generate(context.Background(), core.AgentRequest{Prompt: "hello"})
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if resp.Reply == "" {
		t.Fatalf("empty reply")
	}
}
