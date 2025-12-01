package agent_test

import (
	"context"
	"testing"

	"github.com/joelklabo/buddy/internal/agents/echo"
	"github.com/joelklabo/buddy/internal/core"
)

func TestAgentConformanceSimple(t *testing.T) {
	ag := echo.New()
	resp, err := ag.Generate(context.Background(), core.AgentRequest{Prompt: "ping"})
	if err != nil {
		t.Fatalf("agent %T failed: %v", ag, err)
	}
	if resp.Reply == "" {
		t.Fatalf("agent %T returned empty reply", ag)
	}
}
