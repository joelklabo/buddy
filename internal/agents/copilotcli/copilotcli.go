package copilotcli

// Package copilotcli adapts GitHub Copilot CLI (`gh copilot suggest`) to the Agent interface.

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"nostr-codex-runner/internal/core"
)

// Config controls the Copilot CLI agent.
type Config struct {
	Binary         string
	WorkingDir     string
	TimeoutSeconds int
}

// Agent shells out to `gh copilot suggest -p "<prompt>"`.
type Agent struct {
	cfg Config
}

// New constructs a Copilot CLI agent with sane defaults.
func New(cfg Config) *Agent {
	if cfg.Binary == "" {
		cfg.Binary = "gh"
	}
	if cfg.TimeoutSeconds == 0 {
		cfg.TimeoutSeconds = 120
	}
	return &Agent{cfg: cfg}
}

func (a *Agent) Generate(ctx context.Context, req core.AgentRequest) (core.AgentResponse, error) {
	timeout := time.Duration(a.cfg.TimeoutSeconds) * time.Second
	cctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	args := []string{"copilot", "suggest", "-p", strings.TrimSpace(req.Prompt)}
	cmd := exec.CommandContext(cctx, a.cfg.Binary, args...)
	if a.cfg.WorkingDir != "" {
		cmd.Dir = a.cfg.WorkingDir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if cctx.Err() == context.DeadlineExceeded {
			return core.AgentResponse{}, fmt.Errorf("copilot timeout")
		}
		return core.AgentResponse{}, fmt.Errorf("copilot failed: %v (stderr: %s)", err, stderr.String())
	}

	reply := strings.TrimSpace(stdout.String())
	if reply == "" {
		reply = "(no output from copilot)"
	}

	return core.AgentResponse{
		Reply:       reply,
		ActionCalls: nil,
	}, nil
}
