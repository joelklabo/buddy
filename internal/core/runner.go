package core

import (
    "context"
    "errors"
    "fmt"
    "log/slog"
    "sync"
    "time"
)

// Runner wires transports, agent, and actions together.
type Runner struct {
	transports   []Transport
	transportMap map[string]Transport
	agent        Agent
	actions      map[string]Action
	actionSpecs  []ActionSpec
	logger       *slog.Logger

	reqTimeout    time.Duration
	actionTimeout time.Duration
}

// RunnerOption configures a Runner.
type RunnerOption func(*Runner)

// WithRequestTimeout overrides the per-agent request timeout.
func WithRequestTimeout(d time.Duration) RunnerOption {
	return func(r *Runner) { r.reqTimeout = d }
}

// WithActionTimeout overrides the per-action timeout.
func WithActionTimeout(d time.Duration) RunnerOption {
	return func(r *Runner) { r.actionTimeout = d }
}

// NewRunner constructs a Runner. If logger is nil, slog.Default is used.
func NewRunner(transports []Transport, agent Agent, actions []Action, logger *slog.Logger, opts ...RunnerOption) *Runner {
	if logger == nil {
		logger = slog.Default()
	}

	tmap := make(map[string]Transport, len(transports))
	for _, t := range transports {
		tmap[t.ID()] = t
	}

	amap := make(map[string]Action, len(actions))
	specs := make([]ActionSpec, 0, len(actions))
	for _, a := range actions {
		amap[a.Name()] = a
		specs = append(specs, ActionSpec{
			Name:         a.Name(),
			Capabilities: a.Capabilities(),
		})
	}

	r := &Runner{
		transports:    transports,
		transportMap:  tmap,
		agent:         agent,
		actions:       amap,
		actionSpecs:   specs,
		logger:        logger,
		reqTimeout:    15 * time.Minute,
		actionTimeout: 2 * time.Minute,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// Start launches transports and processes inbound messages until ctx is done.
func (r *Runner) Start(ctx context.Context) error {
	inbound := make(chan InboundMessage, 128)
	var wg sync.WaitGroup
	errCh := make(chan error, len(r.transports))

	for _, t := range r.transports {
		wg.Add(1)
		go func(tr Transport) {
			defer wg.Done()
			if err := tr.Start(ctx, inbound); err != nil {
				errCh <- fmt.Errorf("transport %s: %w", tr.ID(), err)
			}
		}(t)
	}

	// Processor loop
	go func() {
		<-ctx.Done()
		close(inbound)
	}()

	for msg := range inbound {
		r.handleMessage(ctx, msg)
	}

	wg.Wait()

	select {
    case err := <-errCh:
        if errors.Is(err, context.Canceled) {
            return nil
        }
        return err
    default:
        if errors.Is(ctx.Err(), context.Canceled) {
            return nil
        }
        return ctx.Err()
    }
}

func (r *Runner) handleMessage(parent context.Context, msg InboundMessage) {
	log := r.logger.With(
		slog.String("transport", msg.Transport),
		slog.String("sender", msg.Sender),
		slog.String("thread", msg.ThreadID),
	)

	reqCtx := parent
	if r.reqTimeout > 0 {
		var cancel context.CancelFunc
		reqCtx, cancel = context.WithTimeout(parent, r.reqTimeout)
		defer cancel()
	}

	req := AgentRequest{
		Prompt:     msg.Text,
		History:    nil,
		Actions:    r.actionSpecs,
		SenderMeta: msg.Meta,
	}

	start := time.Now()
	resp, err := r.agent.Generate(reqCtx, req)
	if err != nil {
		log.Error("agent error", slog.String("err", err.Error()))
		return
	}
	log.Info("agent reply", slog.Duration("ms", time.Since(start)))

	// Execute actions if any
	var actionResults []string
	for _, call := range resp.ActionCalls {
		act, ok := r.actions[call.Name]
		if !ok {
			log.Warn("unknown action", slog.String("action", call.Name))
			continue
		}
		aCtx := reqCtx
		if r.actionTimeout > 0 {
			var cancel context.CancelFunc
			aCtx, cancel = context.WithTimeout(reqCtx, r.actionTimeout)
			defer cancel()
		}
		aStart := time.Now()
		out, err := act.Invoke(aCtx, call.Args)
		if err != nil {
			log.Error("action error", slog.String("action", call.Name), slog.String("err", err.Error()))
			continue
		}
		log.Info("action ok", slog.String("action", call.Name), slog.Duration("ms", time.Since(aStart)))
		if len(out) > 0 {
			actionResults = append(actionResults, fmt.Sprintf("[%s]\n%s", call.Name, string(out)))
		}
	}

	finalText := resp.Reply
	if len(actionResults) > 0 {
		finalText = finalText + "\n\n" + joinStrings(actionResults, "\n\n")
	}

	outMsg := OutboundMessage{
		Transport: msg.Transport,
		Recipient: msg.Sender,
		Text:      finalText,
		ThreadID:  msg.ThreadID,
	}

	tr, ok := r.transportMap[msg.Transport]
	if !ok {
		log.Error("no transport for outbound", slog.String("transport", msg.Transport))
		return
	}
	if err := tr.Send(reqCtx, outMsg); err != nil {
		log.Error("send error", slog.String("err", err.Error()))
	}
}

func joinStrings(parts []string, sep string) string {
	if len(parts) == 0 {
		return ""
	}
	if len(parts) == 1 {
		return parts[0]
	}
	out := parts[0]
	for _, p := range parts[1:] {
		out += sep + p
	}
	return out
}
