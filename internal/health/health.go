package health

import (
	"context"
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"time"
)

// Start launches a simple /health endpoint. If addr is empty, it is a no-op.
// It returns the actual listening address (useful if addr ends with :0).
func Start(ctx context.Context, addr, version string, logger *slog.Logger) (string, error) {
	if addr == "" {
		return "", nil
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return "", err
	}
	actual := ln.Addr().String()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"version": version,
		})
	})

	srv := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	go func() {
		if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			logger.Error("health server error", slog.String("err", err.Error()))
		}
	}()

	logger.Info("health server listening", slog.String("addr", actual))
	return actual, nil
}
