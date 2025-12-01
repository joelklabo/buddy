package http

import (
	"context"
	"testing"

	"github.com/joelklabo/buddy/internal/core"
)

func TestHTTPStub(t *testing.T) {
	ag := New(Config{APIBase: "https://example.com", Model: "gpt"})
	_, err := ag.Generate(context.Background(), core.AgentRequest{Prompt: "hi"})
	if err == nil {
		t.Fatalf("expected not implemented error")
	}
}
