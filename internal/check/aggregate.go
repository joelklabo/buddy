package check

import "github.com/joelklabo/buddy/internal/config"

// AggregateDeps collects deps for a given config and optional preset name.
// If presetName matches an embedded preset, deps for that preset are included.
func AggregateDeps(cfg *config.Config, presetName string, presetDeps map[string][]config.Dep) []config.Dep {
	var out []config.Dep
	appendDeps := func(m map[string][]config.Dep, key string) {
		if deps, ok := m[key]; ok {
			out = append(out, deps...)
		}
	}

	// From explicit config deps blocks
	appendDeps(cfg.Deps.Transports, transportType(cfg))
	appendDeps(cfg.Deps.Agents, cfg.Agent.Type)
	for _, a := range cfg.Actions {
		appendDeps(cfg.Deps.Actions, a.Type)
	}

	// From preset deps map if provided
	if presetName != "" {
		if deps, ok := presetDeps[presetName]; ok {
			out = append(out, deps...)
		}
	}

	return dedupe(out)
}

func transportType(cfg *config.Config) string {
	if len(cfg.Transports) == 0 {
		return ""
	}
	return cfg.Transports[0].Type
}

func dedupe(deps []config.Dep) []config.Dep {
	seen := map[string]config.Dep{}
	for _, d := range deps {
		key := d.Type + "|" + d.Name
		// later entries override earlier ones
		seen[key] = d
	}
	out := make([]config.Dep, 0, len(seen))
	for _, d := range seen {
		out = append(out, d)
	}
	return out
}
