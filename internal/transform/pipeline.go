package transform

import "fmt"

// Stage is a named transformation step within a Pipeline.
type Stage struct {
	Name        string
	Transformer *Transformer
}

// Pipeline chains multiple Transformer instances, applying them in order.
type Pipeline struct {
	stages []Stage
}

// NewPipeline builds a Pipeline from a slice of named Configs.
// Each entry must have a non-empty Name and a valid Config.
func NewPipeline(entries []struct {
	Name string
	Cfg  Config
}) (*Pipeline, error) {
	p := &Pipeline{}
	for i, e := range entries {
		if e.Name == "" {
			return nil, fmt.Errorf("transform: stage %d has an empty name", i)
		}
		tr, err := New(e.Cfg)
		if err != nil {
			return nil, fmt.Errorf("transform: stage %q: %w", e.Name, err)
		}
		p.stages = append(p.stages, Stage{Name: e.Name, Transformer: tr})
	}
	return p, nil
}

// Apply runs the payload through every stage in order and returns the result.
func (p *Pipeline) Apply(payload string) string {
	for _, s := range p.stages {
		payload = s.Transformer.Apply(payload)
	}
	return payload
}

// Stages returns a copy of the pipeline's stage list (read-only introspection).
func (p *Pipeline) Stages() []Stage {
	out := make([]Stage, len(p.stages))
	copy(out, p.stages)
	return out
}
