// Package mock implements an in-memory transport for tests.
package mock

import (
	"context"
	"github.com/joelklabo/buddy/internal/core"
)

// Transport is an in-memory transport for tests.
type Transport struct {
	id       string
	Inbound  chan core.InboundMessage
	Outbound chan core.OutboundMessage
}

func New(id string) *Transport {
	if id == "" {
		id = "mock"
	}
	return &Transport{
		id:       id,
		Inbound:  make(chan core.InboundMessage, 32),
		Outbound: make(chan core.OutboundMessage, 32),
	}
}

func (t *Transport) ID() string { return t.id }

func (t *Transport) Start(ctx context.Context, in chan<- core.InboundMessage) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-t.Inbound:
			if !ok {
				return nil
			}
			in <- msg
		}
	}
}

func (t *Transport) Send(ctx context.Context, msg core.OutboundMessage) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case t.Outbound <- msg:
		return nil
	}
}
