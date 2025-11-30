package health

import (
	"context"
	"net/http"
	"testing"
	"time"

	"log/slog"
)

func TestHealthServer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	addr, err := Start(ctx, "127.0.0.1:0", "testver", slog.Default())
	if err != nil {
		t.Fatalf("start health: %v", err)
	}
	if addr == "" {
		t.Fatalf("expected addr")
	}

	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://" + addr + "/health")
	if err != nil {
		t.Fatalf("get health: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
	cancel()
}
