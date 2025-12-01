package check

import (
	"testing"

	"github.com/joelklabo/buddy/internal/config"
)

func TestAggregateDeps(t *testing.T) {
	cfg := &config.Config{
		Deps: config.DepSet{
			Transports: map[string][]config.Dep{
				"nostr": {{Name: "curl", Type: "binary"}},
			},
			Agents: map[string][]config.Dep{
				"http": {{Name: "openssl", Type: "binary"}},
			},
			Actions: map[string][]config.Dep{
				"shell": {{Name: "bash", Type: "binary"}},
			},
		},
		Transports: []config.TransportConfig{{Type: "nostr"}},
		Agent:      config.AgentConfig{Type: "http"},
		Actions:    []config.ActionConfig{{Type: "shell"}},
	}
	presetDeps := map[string][]config.Dep{
		"mock-echo": {{Name: "none", Type: "binary"}},
	}
	out := AggregateDeps(cfg, "mock-echo", presetDeps)
	if len(out) != 4 {
		t.Fatalf("expected 4 deps, got %d", len(out))
	}
}
