package mailgun

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/joelklabo/buddy/internal/core"
)

func TestHandlerAcceptsValidRequest(t *testing.T) {
	cfg := Config{
		Domain:       "mg.example.com",
		APIKey:       "test-api",
		SigningKey:   "test-key",
		AllowSenders: []string{"alice@example.com"},
		Listen:       ":0",
		Path:         "/email/inbound",
	}
	cfg.Defaults()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}
	tr, err := New(cfg)
	if err != nil {
		t.Fatalf("new transport: %v", err)
	}

	inbound := make(chan core.InboundMessage, 1)
	handler := tr.Handler(inbound)

	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	form := url.Values{}
	timestamp := "1700000000"
	token := "abcdef"
	sig := hmacHex(timestamp, token, cfg.SigningKey)

	form.Set("timestamp", timestamp)
	form.Set("token", token)
	form.Set("signature", sig)
	form.Set("sender", "alice@example.com")
	form.Set("stripped-text", "hello")
	form.Set("Message-Id", "<m1>")
	form.Set("subject", "hi")

	resp, err := http.PostForm(ts.URL, form)
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", resp.StatusCode)
	}

	select {
	case msg := <-inbound:
		if msg.Text != "hello" {
			t.Fatalf("unexpected text: %q", msg.Text)
		}
		if msg.Sender != "alice@example.com" {
			t.Fatalf("unexpected sender: %s", msg.Sender)
		}
	case <-time.After(time.Second):
		t.Fatalf("no inbound message received")
	}
}

func TestHandlerRejectsBadSignature(t *testing.T) {
	cfg := Config{
		Domain:       "mg.example.com",
		APIKey:       "test-api",
		SigningKey:   "test-key",
		AllowSenders: []string{"alice@example.com"},
		Listen:       ":0",
		Path:         "/email/inbound",
	}
	cfg.Defaults()
	tr, _ := New(cfg)

	inbound := make(chan core.InboundMessage, 1)
	ts := httptest.NewServer(tr.Handler(inbound))
	t.Cleanup(ts.Close)

	form := url.Values{}
	form.Set("timestamp", "1700000000")
	form.Set("token", "abcdef")
	form.Set("signature", "deadbeef")
	form.Set("sender", "alice@example.com")

	resp, err := http.PostForm(ts.URL, form)
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func hmacHex(ts, token, key string) string {
	return hexEncode(ts, token, key)
}

// hexEncode uses verifySignature path to avoid duplicating logic.
func hexEncode(timestamp, token, key string) string {
	// reuse verifySignature internals indirectly by constructing expected
	mac := hmacSha(timestamp, token, key)
	return mac
}

func hmacSha(timestamp, token, key string) string {
	// direct helper used by test; matches verifySignature logic
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(timestamp))
	mac.Write([]byte(token))
	return hex.EncodeToString(mac.Sum(nil))
}
