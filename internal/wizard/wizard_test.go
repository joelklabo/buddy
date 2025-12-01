package wizard

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunWritesConfig(t *testing.T) {
	td := t.TempDir()
	path := filepath.Join(td, "config.yaml")
	p := &StubPrompter{
		Inputs:    []string{"wss://relay.example", "deadbeef"},
		Passwords: []string{"abcd1234"},
		Selects:   []string{"echo"},
		Confirms:  []bool{true},
	}
	got, err := Run(context.Background(), path, p)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if got != path {
		t.Fatalf("expected path %s, got %s", path, got)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	content := string(data)
	if !containsAll(content, []string{"private_key: abcd1234", "allowed_pubkeys:", "shell"}) {
		t.Fatalf("config missing expected fields:\n%s", content)
	}
}

func TestRunRequiresAllowedPubkey(t *testing.T) {
	td := t.TempDir()
	path := filepath.Join(td, "config.yaml")
	p := &StubPrompter{
		Inputs:    []string{"wss://relay.example", ""},
		Passwords: []string{"abcd1234"},
		Selects:   []string{"echo"},
		Confirms:  []bool{false},
	}
	_, err := Run(context.Background(), path, p)
	if err == nil {
		t.Fatalf("expected error for missing allowed pubkeys")
	}
}

func containsAll(s string, needles []string) bool {
	for _, n := range needles {
		if !strings.Contains(s, n) {
			return false
		}
	}
	return true
}
