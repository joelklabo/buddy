// Package whatsapp implements a transport backed by the WhatsApp Cloud API (Meta Graph).
package whatsapp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"nostr-codex-runner/internal/core"
	transport "nostr-codex-runner/internal/transports"
)

// Config holds WhatsApp Cloud API settings.
// Access tokens should be long-lived app tokens or user tokens with whatsapp_business_messaging scope.
type Config struct {
	ID             string   `json:"id"`
	PhoneNumberID  string   `json:"phone_number_id"`
	AccessToken    string   `json:"access_token"`
	VerifyToken    string   `json:"verify_token"`
	Listen         string   `json:"listen"`          // e.g. ":8082"
	BaseURL        string   `json:"base_url"`        // default: https://graph.facebook.com/v18.0
	AllowedNumbers []string `json:"allowed_numbers"` // optional allowlist of sender phone numbers (E.164 without +)
}

// Transport implements the core.Transport interface for WhatsApp.
type Transport struct {
	cfg      Config
	client   *http.Client
	log      *slog.Logger
	srv      *http.Server
	started  chan string // publishes listen addr for tests/observability
	startMu  sync.Mutex
	shutdown func(context.Context) error
}

// New constructs a WhatsApp transport.
func New(cfg Config, logger *slog.Logger) (*Transport, error) {
	if cfg.ID == "" {
		cfg.ID = "whatsapp"
	}
	if cfg.PhoneNumberID == "" {
		return nil, errors.New("whatsapp: phone_number_id required")
	}
	if cfg.AccessToken == "" {
		return nil, errors.New("whatsapp: access_token required")
	}
	if cfg.VerifyToken == "" {
		return nil, errors.New("whatsapp: verify_token required")
	}
	if cfg.Listen == "" {
		cfg.Listen = ":8082"
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://graph.facebook.com/v18.0"
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &Transport{
		cfg:     cfg,
		client:  &http.Client{Timeout: 15 * time.Second},
		log:     logger.With("transport", "whatsapp"),
		started: make(chan string, 1),
	}, nil
}

// ID returns the transport identifier.
func (t *Transport) ID() string { return t.cfg.ID }

// Addr returns the bound listen address (useful in tests).
func (t *Transport) Addr() string {
	select {
	case addr := <-t.started:
		t.started <- addr
		return addr
	default:
		return ""
	}
}

// Start launches the webhook listener and blocks until ctx is canceled or a fatal error occurs.
func (t *Transport) Start(ctx context.Context, inbound chan<- core.InboundMessage) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Verification callback
			mode := r.URL.Query().Get("hub.mode")
			token := r.URL.Query().Get("hub.verify_token")
			challenge := r.URL.Query().Get("hub.challenge")
			if mode == "subscribe" && token == t.cfg.VerifyToken {
				w.WriteHeader(http.StatusOK)
				_, _ = io.WriteString(w, challenge)
				return
			}
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		case http.MethodPost:
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}
			if err := t.handleWebhook(ctx, inbound, body); err != nil {
				t.log.Warn("webhook handling failed", "err", err)
			}
			w.WriteHeader(http.StatusOK)
			return
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	ln, err := net.Listen("tcp", t.cfg.Listen)
	if err != nil {
		return fmt.Errorf("whatsapp listen: %w", err)
	}
	addr := ln.Addr().String()
	t.started <- addr

	srv := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	t.startMu.Lock()
	t.srv = srv
	t.shutdown = srv.Shutdown
	t.startMu.Unlock()

	errCh := make(chan error, 1)
	go func() {
		if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
		return nil
	case err := <-errCh:
		return err
	}
}

// Send posts a message via the WhatsApp Cloud API.
func (t *Transport) Send(ctx context.Context, msg core.OutboundMessage) error {
	payload := map[string]any{
		"messaging_product": "whatsapp",
		"to":                msg.Recipient,
		"type":              "text",
		"text": map[string]any{
			"body": msg.Text,
		},
	}
	if msg.ThreadID != "" {
		payload["context"] = map[string]any{"message_id": msg.ThreadID}
	}
	body, _ := json.Marshal(payload)
	url := fmt.Sprintf("%s/%s/messages", strings.TrimRight(t.cfg.BaseURL, "/"), t.cfg.PhoneNumberID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+t.cfg.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("whatsapp send failed: %s", strings.TrimSpace(string(b)))
	}
	return nil
}

// handleWebhook parses incoming webhook payloads and emits InboundMessage events.
func (t *Transport) handleWebhook(ctx context.Context, inbound chan<- core.InboundMessage, raw []byte) error {
	var webhook struct {
		Entry []struct {
			Changes []struct {
				Value struct {
					Messages []struct {
						From string `json:"from"`
						ID   string `json:"id"`
						Text struct {
							Body string `json:"body"`
						} `json:"text"`
						Context struct {
							ID string `json:"id"`
						} `json:"context"`
					} `json:"messages"`
				} `json:"value"`
			} `json:"changes"`
		} `json:"entry"`
	}
	if err := json.Unmarshal(raw, &webhook); err != nil {
		return fmt.Errorf("decode webhook: %w", err)
	}
	for _, entry := range webhook.Entry {
		for _, change := range entry.Changes {
			for _, m := range change.Value.Messages {
				if len(t.cfg.AllowedNumbers) > 0 && !contains(t.cfg.AllowedNumbers, m.From) {
					t.log.Warn("rejecting sender not in allowlist", "from", m.From)
					continue
				}
				im := core.InboundMessage{
					Transport: t.ID(),
					Sender:    m.From,
					Text:      m.Text.Body,
					ThreadID:  m.Context.ID,
				}
				select {
				case inbound <- im:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}
	}
	return nil
}

func contains(list []string, v string) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}

func init() {
	transport.MustRegister("whatsapp", func(cfg any) (core.Transport, error) {
		c, ok := cfg.(Config)
		if !ok {
			return nil, fmt.Errorf("whatsapp: invalid config type %T", cfg)
		}
		return New(c, nil)
	})
}
