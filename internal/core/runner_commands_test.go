package core

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
	"time"

	"nostr-codex-runner/internal/store"
)

// mock action for shell
type shellAction struct{ invoked bool }

func (s *shellAction) Name() string           { return "shell" }
func (s *shellAction) Capabilities() []string { return []string{"shell:exec"} }
func (s *shellAction) Help() string           { return "/shell <cmd> — mock shell" }
func (s *shellAction) Invoke(ctx context.Context, args json.RawMessage) (json.RawMessage, error) {
	s.invoked = true
	return json.RawMessage(`"ok"`), nil
}

// mock store for status/use/new
type memoryStore struct{ active map[string]store.SessionState }

func (m *memoryStore) Active(sender string) (store.SessionState, bool, error) {
	if st, ok := m.active[sender]; ok {
		return st, true, nil
	}
	return store.SessionState{}, false, nil
}
func (m *memoryStore) SaveActive(sender, sessionID string) error {
	if m.active == nil {
		m.active = map[string]store.SessionState{}
	}
	m.active[sender] = store.SessionState{SessionID: sessionID, UpdatedAt: time.Now()}
	return nil
}
func (m *memoryStore) ClearActive(sender string) error {
	delete(m.active, sender)
	return nil
}

// satisfy AuditLogger
func (m *memoryStore) AppendAudit(action, sender, outcome string, dur time.Duration) error {
	return nil
}

// extra methods to satisfy StoreAPI (no-ops for tests)
func (m *memoryStore) LastCursor(pubkey string) (time.Time, error)   { return time.Time{}, nil }
func (m *memoryStore) SaveCursor(pubkey string, ts time.Time) error  { return nil }
func (m *memoryStore) AlreadyProcessed(eventID string) (bool, error) { return false, nil }
func (m *memoryStore) MarkProcessed(eventID string) error            { return nil }
func (m *memoryStore) RecentMessageSeen(pubkey, message string, window time.Duration) (bool, error) {
	return false, nil
}

func TestRunnerHelpIncludesActionHelp(t *testing.T) {
	act := &shellAction{}
	r := NewRunner(nil, &mockAgent{reply: "hi"}, []Action{act}, slog.Default())
	msg := InboundMessage{Transport: "mock", Sender: "alice", Text: "/help", ThreadID: "t1"}
	out := captureSend(r, msg)
	if !contains(out, "/shell <cmd>") {
		t.Fatalf("expected shell help in %q", out)
	}
}

func TestRunnerShellCommandInvokesAction(t *testing.T) {
	act := &shellAction{}
	r := NewRunner([]Transport{&mockTransport{id: "mock"}}, &mockAgent{reply: "hi"}, []Action{act}, slog.Default())
	msg := InboundMessage{Transport: "mock", Sender: "alice", Text: "/shell ls", ThreadID: "t1"}
	_ = captureSend(r, msg)
	if !act.invoked {
		t.Fatalf("shell action not invoked")
	}
}

func TestRunnerStatusAndUse(t *testing.T) {
	st := &memoryStore{active: map[string]store.SessionState{"alice": {SessionID: "sess1", UpdatedAt: time.Now()}}}
	r := NewRunner(nil, &mockAgent{reply: "hi"}, nil, slog.Default(), WithStore(st))
	msg := InboundMessage{Transport: "mock", Sender: "alice", Text: "/status", ThreadID: "t1"}
	out := captureSend(r, msg)
	if !contains(out, "sess1") {
		t.Fatalf("expected status to mention sess1, got %q", out)
	}
	msg2 := InboundMessage{Transport: "mock", Sender: "alice", Text: "/use sess2", ThreadID: "t1"}
	_ = captureSend(r, msg2)
	if st.active["alice"].SessionID != "sess2" {
		t.Fatalf("use did not update active")
	}
}

func TestRunnerNewClearsSession(t *testing.T) {
	st := &memoryStore{active: map[string]store.SessionState{"alice": {SessionID: "sess1", UpdatedAt: time.Now()}}}
	r := NewRunner(nil, &mockAgent{reply: "hi"}, nil, slog.Default(), WithStore(st))
	msg := InboundMessage{Transport: "mock", Sender: "alice", Text: "/new", ThreadID: "t1"}
	_ = captureSend(r, msg)
	if _, ok := st.active["alice"]; ok {
		t.Fatalf("expected active cleared")
	}
}

func contains(s, sub string) bool { return strings.Contains(s, sub) }

// captureSend runs handleCommand path and returns last sent text
func captureSend(r *Runner, msg InboundMessage) string {
	outCh := make(chan OutboundMessage, 1)
	r.transportMap = map[string]Transport{"mock": &transportSpy{out: outCh}}
	r.handleCommand(context.Background(), msg, slog.Default())
	select {
	case out := <-outCh:
		return out.Text
	default:
		return ""
	}
}

type transportSpy struct {
	out chan<- OutboundMessage
}

func (t *transportSpy) ID() string { return "mock" }
func (t *transportSpy) Start(ctx context.Context, inbound chan<- InboundMessage) error {
	return nil
}
func (t *transportSpy) Send(ctx context.Context, msg OutboundMessage) error {
	t.out <- msg
	return nil
}
