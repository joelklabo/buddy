package health

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"testing"
	"time"
)

func TestHealthEndpointServes(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	addr, err := Start(ctx, "127.0.0.1:18101", "vtest", logger)
	if err != nil {
		t.Fatalf("start health: %v", err)
	}

	// wait briefly for server
	time.Sleep(50 * time.Millisecond)
	resp, err := http.Get("http://" + addr + "/health")
	if err != nil {
		t.Fatalf("http get: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("unexpected status %d", resp.StatusCode)
	}
	_ = resp.Body.Close()
}
