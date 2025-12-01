package whatsapp

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"nostr-codex-runner/internal/core"
)

func TestWebhookInbound(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tr, err := New(Config{
		ID:            "wa",
		PhoneNumberID: "12345",
		AccessToken:   "t",
		VerifyToken:   "verify",
		Listen:        "127.0.0.1:0",
		AllowedNumbers: []string{
			"15555550100",
		},
	}, nil)
	if err != nil {
		t.Fatalf("new: %v", err)
	}

	inbound := make(chan core.InboundMessage, 1)
	go func() {
		if err := tr.Start(ctx, inbound); err != nil && ctx.Err() == nil {
			t.Errorf("start: %v", err)
		}
	}()

	// wait for server to bind
	addr := ""
	for i := 0; i < 50 && addr == ""; i++ {
		time.Sleep(10 * time.Millisecond)
		addr = tr.Addr()
	}
	if addr == "" {
		t.Fatalf("server did not start")
	}
	url := "http://" + addr + "/"

	// verify challenge
	resp, err := http.Get(url + "?hub.mode=subscribe&hub.verify_token=verify&hub.challenge=42")
	if err != nil {
		t.Fatalf("verify get: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	// send webhook payload
	payload := map[string]any{
		"entry": []any{
			map[string]any{
				"changes": []any{
					map[string]any{
						"value": map[string]any{
							"messages": []any{
								map[string]any{
									"from": "15555550100",
									"id":   "wamid.ABCD",
									"text": map[string]any{"body": "hello wa"},
								},
							},
						},
					},
				},
			},
		},
	}
	b, _ := json.Marshal(payload)
	r, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 post, got %d", r.StatusCode)
	}

	select {
	case msg := <-inbound:
		if msg.Sender != "15555550100" || msg.Text != "hello wa" {
			t.Fatalf("unexpected inbound: %+v", msg)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("no inbound message received")
	}
}

func TestSendOutbound(t *testing.T) {
	var got struct {
		To   string `json:"to"`
		Body string `json:"body"`
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer tok" {
			t.Fatalf("auth header missing")
		}
		if r.Method != http.MethodPost {
			t.Fatalf("method %s", r.Method)
		}
		var payload struct {
			To   string `json:"to"`
			Text struct {
				Body string `json:"body"`
			} `json:"text"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode: %v", err)
		}
		got.To = payload.To
		got.Body = payload.Text.Body
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	tr, err := New(Config{
		PhoneNumberID: "pnid",
		AccessToken:   "tok",
		VerifyToken:   "v",
		BaseURL:       srv.URL,
	}, nil)
	if err != nil {
		t.Fatalf("new: %v", err)
	}

	if err := tr.Send(context.Background(), core.OutboundMessage{
		Recipient: "1555",
		Text:      "hi there",
	}); err != nil {
		t.Fatalf("send: %v", err)
	}
	if got.To != "1555" || got.Body != "hi there" {
		t.Fatalf("unexpected send payload: %+v", got)
	}
}
