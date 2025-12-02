package core_test

import (
	"context"
	"testing"
	"time"

	"github.com/joelklabo/buddy/internal/agents/echo"
	"github.com/joelklabo/buddy/internal/core"
	tmock "github.com/joelklabo/buddy/internal/transports/mock"
)

// Load sanity: ensure runner handles a burst of messages without blocking indefinitely.
func TestRunnerHandlesBurst(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	transport := tmock.New("mock-load")
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

	// Burst 50 messages.
	for i := 0; i < 50; i++ {
		transport.Inbound <- core.InboundMessage{
			Transport: transport.ID(),
			Sender:    "alice",
			Text:      "/new msg",
			ThreadID:  "t1",
		}
	}

	// Expect same number of replies within a timeout.
	received := 0
	timeout := time.After(3 * time.Second)
	for received < 50 {
		select {
		case <-transport.Outbound:
			received++
		case <-timeout:
			t.Fatalf("timed out; received %d of 50", received)
		}
	}

	cancel()
	<-done
}
