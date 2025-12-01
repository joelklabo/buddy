package codexcli

import (
	"context"
	"errors"
	"testing"

	"github.com/joelklabo/buddy/internal/codex"
	"github.com/joelklabo/buddy/internal/core"
)

type fakeRunner struct {
	reply string
	err   error
}

func (f *fakeRunner) Run(ctx context.Context, sessionID string, prompt string) (codex.Result, error) {
	return codex.Result{Reply: f.reply, SessionID: "s1"}, f.err
}

func (f *fakeRunner) ContextWithTimeout(parent context.Context) (context.Context, context.CancelFunc) {
	return context.WithCancel(parent)
}

func TestGenerateErrorOnEmptyPrompt(t *testing.T) {
	ag := &Agent{runner: &fakeRunner{}}
	if _, err := ag.Generate(context.Background(), core.AgentRequest{}); err == nil {
		t.Fatalf("expected error")
	}
}

func TestGeneratePassesPrompt(t *testing.T) {
	ag := &Agent{runner: &fakeRunner{reply: "ok"}}
	out, err := ag.Generate(context.Background(), core.AgentRequest{Prompt: "hi"})
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if out.Reply != "ok" || out.SessionID != "s1" {
		t.Fatalf("unexpected reply %+v", out)
	}
}

func TestGenerateRunnerError(t *testing.T) {
	ag := &Agent{runner: &fakeRunner{err: errors.New("boom")}}
	if _, err := ag.Generate(context.Background(), core.AgentRequest{Prompt: "hi"}); err == nil {
		t.Fatalf("expected error")
	}
}
