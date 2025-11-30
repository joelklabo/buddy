package core_test

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"strings"
	"testing"
	"time"

	"nostr-codex-runner/internal/actions/shell"
	"nostr-codex-runner/internal/core"
	tmock "nostr-codex-runner/internal/transports/mock"
)

type stubAgent struct {
	resp core.AgentResponse
}

func (s *stubAgent) Generate(ctx context.Context, req core.AgentRequest) (core.AgentResponse, error) {
	return s.resp, nil
}

func TestRunnerWithMockTransportAndShellAction(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tr := tmock.New("mock")
	ag := &stubAgent{resp: core.AgentResponse{
		Reply: "base",
		ActionCalls: []core.ActionCall{
			{Name: "shell", Args: json.RawMessage(`{"command":"echo hi"}`)},
		},
	}}
	sh := shell.New(shell.Config{Allowed: []string{"echo "}, TimeoutSeconds: 5})

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	r := core.NewRunner([]core.Transport{tr}, ag, []core.Action{sh}, logger, core.WithActionTimeout(2*time.Second))

	done := make(chan struct{})
	go func() {
		_ = r.Start(ctx)
		close(done)
	}()

	tr.Inbound <- core.InboundMessage{Transport: "mock", Sender: "alice", Text: "please run shell", ThreadID: "t1"}

	var out core.OutboundMessage
	select {
	case out = <-tr.Outbound:
	case <-time.After(3 * time.Second):
		t.Fatal("no outbound message")
	}

	cancel()
	<-done

	if !strings.Contains(out.Text, "hi") {
		t.Fatalf("expected shell output, got %q", out.Text)
	}
}
