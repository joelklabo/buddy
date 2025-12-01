package metrics

import (
	"context"
	"log/slog"
	"net/http"
	"testing"
	"time"
)

func TestStartNoListen(t *testing.T) {
	if err := Start(context.Background(), "", slog.Default()); err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}
}

func TestStartAndExpose(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	addr := "127.0.0.1:18099"
	if err := Start(ctx, addr, slog.Default()); err != nil {
		t.Fatalf("start metrics: %v", err)
	}
	var resp *http.Response
	var err error
	for i := 0; i < 10; i++ {
		resp, err = http.Get("http://" + addr + "/metrics")
		if err == nil {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	if err != nil {
		t.Fatalf("get metrics: %v", err)
	}
	if cerr := resp.Body.Close(); cerr != nil {
		t.Fatalf("close body: %v", cerr)
	}
}
