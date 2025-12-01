package wizard

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"gopkg.in/yaml.v3"

	"github.com/joelklabo/buddy/internal/config"
)

// Prompter abstracts survey for testability.
type Prompter interface {
	AskSelect(label string, options []string, def string) (string, error)
	AskInput(label, def string) (string, error)
	AskPassword(label string) (string, error)
	AskConfirm(label string, def bool) (bool, error)
}

// Run executes the interactive wizard and writes a config file.
func Run(ctx context.Context, path string, p Prompter) (string, error) {
	_ = ctx // reserved for future use (cancellation)
	if p == nil {
		p = &surveyPrompter{}
	}

	cfgPath, err := resolveConfigPath(path)
	if err != nil {
		return "", err
	}

	if fileExists(cfgPath) {
		overwrite, err := p.AskConfirm(fmt.Sprintf("%s exists. Overwrite?", cfgPath), false)
		if err != nil {
			return "", err
		}
		if !overwrite {
			return "", fmt.Errorf("aborted: config exists at %s", cfgPath)
		}
	}

	reg := GetRegistry()

	transportOptions := transportNames(reg.Transports)
	transportChoice, err := p.AskSelect("Transport", transportOptions, defaultChoice("nostr", transportOptions))
	if err != nil {
		return "", err
	}

	var relays string
	var priv string
	var allowedKeys []string

	if transportChoice == "nostr" {
		if relays, err = p.AskInput("Relays (comma-separated)", "wss://relay.damus.io,wss://nos.lol"); err != nil {
			return "", err
		}

		for retries := 0; retries < 2; retries++ {
			priv, err = p.AskPassword("Nostr private key (hex, not nsec)")
			if err != nil {
				return "", err
			}
			if len(priv) > 0 {
				break
			}
		}
		if len(priv) == 0 {
			return "", errors.New("private key is required")
		}

		var allowed string
		for retries := 0; retries < 2; retries++ {
			allowed, err = p.AskInput("Allowed pubkeys (comma-separated hex)", "")
			if err != nil {
				return "", err
			}
			allowedKeys = splitCSV(allowed)
			if len(allowedKeys) > 0 {
				break
			}
		}
		if len(allowedKeys) == 0 {
			return "", errors.New("at least one allowed pubkey required")
		}
	}

	agentOptions := agentNames(reg.Agents)
	agentChoice, err := p.AskSelect("Agent", agentOptions, defaultChoice("http", agentOptions))
	if err != nil {
		return "", err
	}

	enableShell := actionDefault(reg.Actions, "shell")
	if _, err := actionOption(reg.Actions, "shell"); err == nil {
		enableShell, err = p.AskConfirm("Enable shell action? (high risk; trusted operators only)", enableShell)
		if err != nil {
			return "", err
		}
	}
	dryRun, err := p.AskConfirm("Dry-run only (preview config without writing)?", false)
	if err != nil {
		return "", err
	}

	cfg := &config.Config{
		Runner: config.RunnerConfig{
			PrivateKey:     priv,
			AllowedPubkeys: allowedKeys,
		},
		Storage:    config.StorageConfig{Path: defaultStatePath()},
		Transports: []config.TransportConfig{},
		Agent:      config.AgentConfig{Type: agentChoice},
		Actions: []config.ActionConfig{{
			Type:  "readfile",
			Name:  "readfile",
			Roots: []string{"."},
		}},
		Projects: []config.Project{{ID: "default", Name: "default", Path: "."}},
	}

	if transportChoice == "nostr" {
		cfg.Transports = append(cfg.Transports, config.TransportConfig{
			Type:           "nostr",
			ID:             "nostr",
			Relays:         splitCSV(relays),
			PrivateKey:     priv,
			AllowedPubkeys: allowedKeys,
		})
	} else {
		cfg.Transports = append(cfg.Transports, config.TransportConfig{
			Type: "mock",
			ID:   "mock",
		})
	}

	if enableShell {
		cfg.Actions = append(cfg.Actions, config.ActionConfig{Type: "shell", Name: "shell", Workdir: ".", TimeoutSecs: 30, MaxOutput: 4000})
	}

	if cfg.Logging.Level == "" {
		cfg.Logging.Level = "info"
	}
	if cfg.Logging.Format == "" {
		cfg.Logging.Format = "text"
	}

	if dryRun {
		fmt.Printf("Dry run: config NOT written. Target path would be %s\n", cfgPath)
		return cfgPath, nil
	}

	if err := writeConfig(cfgPath, cfg); err != nil {
		return "", err
	}

	return cfgPath, nil
}

func resolveConfigPath(path string) (string, error) {
	if path != "" {
		return path, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "buddy", "config.yaml"), nil
}

func writeConfig(path string, cfg *config.Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("make config dir: %w", err)
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func defaultStatePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "state.db"
	}
	return filepath.Join(home, ".local", "share", "buddy", "state.db")
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// surveyPrompter is the real interactive implementation.
type surveyPrompter struct{}

func (surveyPrompter) AskSelect(label string, options []string, def string) (string, error) {
	sel := def
	prompt := &survey.Select{Message: label, Options: options, Default: def}
	if err := survey.AskOne(prompt, &sel); err != nil {
		return "", err
	}
	return sel, nil
}

func (surveyPrompter) AskInput(label, def string) (string, error) {
	ans := def
	prompt := &survey.Input{Message: label, Default: def}
	if err := survey.AskOne(prompt, &ans); err != nil {
		return "", err
	}
	return ans, nil
}

func (surveyPrompter) AskPassword(label string) (string, error) {
	var ans string
	prompt := &survey.Password{Message: label}
	if err := survey.AskOne(prompt, &ans); err != nil {
		return "", err
	}
	return ans, nil
}

func (surveyPrompter) AskConfirm(label string, def bool) (bool, error) {
	ans := def
	prompt := &survey.Confirm{Message: label, Default: def}
	if err := survey.AskOne(prompt, &ans); err != nil {
		return false, err
	}
	return ans, nil
}

func transportNames(opts []TransportOption) []string {
	names := make([]string, 0, len(opts))
	for _, o := range opts {
		names = append(names, o.Name)
	}
	return names
}

func agentNames(opts []AgentOption) []string {
	names := make([]string, 0, len(opts))
	for _, o := range opts {
		names = append(names, o.Name)
	}
	return names
}

func defaultChoice(defaultVal string, options []string) string {
	for _, opt := range options {
		if opt == defaultVal {
			return defaultVal
		}
	}
	if len(options) > 0 {
		return options[0]
	}
	return defaultVal
}

func actionOption(actions []ActionOption, name string) (ActionOption, error) {
	for _, a := range actions {
		if a.Name == name {
			return a, nil
		}
	}
	return ActionOption{}, fmt.Errorf("action option %s not found", name)
}

func actionDefault(actions []ActionOption, name string) bool {
	a, err := actionOption(actions, name)
	if err != nil {
		return false
	}
	return a.DefaultEnable
}

// StubPrompter is used in tests.
type StubPrompter struct {
	Selects   []string
	Inputs    []string
	Passwords []string
	Confirms  []bool
}

func (s *StubPrompter) popSelect(def string) string {
	if len(s.Selects) == 0 {
		return def
	}
	v := s.Selects[0]
	s.Selects = s.Selects[1:]
	return v
}

func (s *StubPrompter) popInput(def string) string {
	if len(s.Inputs) == 0 {
		return def
	}
	v := s.Inputs[0]
	s.Inputs = s.Inputs[1:]
	return v
}

func (s *StubPrompter) popPassword() string {
	if len(s.Passwords) == 0 {
		return ""
	}
	v := s.Passwords[0]
	s.Passwords = s.Passwords[1:]
	return v
}

func (s *StubPrompter) popConfirm(def bool) bool {
	if len(s.Confirms) == 0 {
		return def
	}
	v := s.Confirms[0]
	s.Confirms = s.Confirms[1:]
	return v
}

func (s *StubPrompter) AskSelect(label string, options []string, def string) (string, error) {
	return s.popSelect(def), nil
}
func (s *StubPrompter) AskInput(label, def string) (string, error) {
	return s.popInput(def), nil
}
func (s *StubPrompter) AskPassword(label string) (string, error) {
	return s.popPassword(), nil
}
func (s *StubPrompter) AskConfirm(label string, def bool) (bool, error) {
	return s.popConfirm(def), nil
}
