package core_test

import (
	"context"
	"testing"
	"time"

	"github.com/joelklabo/buddy/internal/agents/echo"
	"github.com/joelklabo/buddy/internal/core"
	tmock "github.com/joelklabo/buddy/internal/transports/mock"
)

// End-to-end smoke: /new with mock transport + echo agent.
func TestE2E_MockEcho_NewSession(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	transport := tmock.New("mock-e2e")
	agent := echo.New()

	r := core.NewRunner([]core.Transport{transport}, agent, nil, nil,
		core.WithAllowedSenders([]string{"alice"}),
		core.WithMaxReplyChars(2000),
	)

	done := make(chan struct{})
	go func() {
		_ = r.Start(ctx)
		close(done)
	}()

	transport.Inbound <- core.InboundMessage{
		Transport: transport.ID(),
		Sender:    "alice",
		Text:      "/new hello world",
		ThreadID:  "t1",
	}

	select {
	case out := <-transport.Outbound:
		if out.Recipient != "alice" {
			t.Fatalf("expected recipient alice, got %q", out.Recipient)
		}
		if out.Text == "" {
			t.Fatalf("expected reply text")
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timed out waiting for outbound message")
	}

	cancel()
	<-done
}
