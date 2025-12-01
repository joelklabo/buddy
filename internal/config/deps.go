package config

// Dep describes a single prerequisite for an extension.
type Dep struct {
	Name        string `yaml:"name" json:"name"`
	Type        string `yaml:"type" json:"type"`                           // binary|env|url|file|port
	Version     string `yaml:"version,omitempty" json:"version,omitempty"` // semver or range
	Hint        string `yaml:"hint,omitempty" json:"hint,omitempty"`
	Optional    bool   `yaml:"optional,omitempty" json:"optional,omitempty"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
}

// DepSet groups deps for transports, agents, actions, and presets.
type DepSet struct {
	Transports map[string][]Dep `yaml:"transports,omitempty" json:"transports,omitempty"`
	Agents     map[string][]Dep `yaml:"agents,omitempty" json:"agents,omitempty"`
	Actions    map[string][]Dep `yaml:"actions,omitempty" json:"actions,omitempty"`
	Presets    map[string][]Dep `yaml:"presets,omitempty" json:"presets,omitempty"`
}

// MergeDepSets merges multiple DepSets, later entries overriding earlier ones for the same key.
func MergeDepSets(sets ...DepSet) DepSet {
	out := DepSet{
		Transports: map[string][]Dep{},
		Agents:     map[string][]Dep{},
		Actions:    map[string][]Dep{},
		Presets:    map[string][]Dep{},
	}
	for _, s := range sets {
		merge := func(dst map[string][]Dep, src map[string][]Dep) {
			for k, v := range src {
				dst[k] = v
			}
		}
		merge(out.Transports, s.Transports)
		merge(out.Agents, s.Agents)
		merge(out.Actions, s.Actions)
		merge(out.Presets, s.Presets)
	}
	return out
}
