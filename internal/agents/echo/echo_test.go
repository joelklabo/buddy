package echo

import (
	"context"
	"testing"

	"github.com/joelklabo/buddy/internal/core"
)

func TestEchoGenerate(t *testing.T) {
	ag := New()
	out, err := ag.Generate(context.Background(), core.AgentRequest{Prompt: "hi"})
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if out.Reply != "hi" {
		t.Fatalf("unexpected reply %s", out.Reply)
	}
}
