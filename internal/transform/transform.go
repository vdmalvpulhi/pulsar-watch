// Package transform provides message payload transformation utilities
// for pulsar-watch, supporting field masking, truncation, and key remapping.
package transform

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Config holds the transformation rules applied to each message payload.
type Config struct {
	// MaskFields lists JSON field names whose values will be replaced with "***".
	MaskFields []string `yaml:"mask_fields"`
	// TruncateAt is the maximum byte length for the payload string (0 = unlimited).
	TruncateAt int `yaml:"truncate_at"`
	// RenameFields maps original JSON field names to new names.
	RenameFields map[string]string `yaml:"rename_fields"`
}

// Transformer applies a set of transformations to raw message payloads.
type Transformer struct {
	cfg Config
}

// New creates a Transformer from the provided Config.
// Returns an error if the configuration is invalid.
func New(cfg Config) (*Transformer, error) {
	if cfg.TruncateAt < 0 {
		return nil, fmt.Errorf("transform: truncate_at must be >= 0, got %d", cfg.TruncateAt)
	}
	return &Transformer{cfg: cfg}, nil
}

// Apply transforms the given payload string according to the configured rules.
// Non-JSON payloads are only subject to truncation.
func (t *Transformer) Apply(payload string) string {
	result := payload

	if len(t.cfg.MaskFields) > 0 || len(t.cfg.RenameFields) > 0 {
		result = t.transformJSON(result)
	}

	if t.cfg.TruncateAt > 0 && len(result) > t.cfg.TruncateAt {
		result = result[:t.cfg.TruncateAt] + "…"
	}

	return result
}

// transformJSON attempts to parse payload as JSON and applies mask/rename rules.
// If parsing fails the original payload is returned unchanged.
func (t *Transformer) transformJSON(payload string) string {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &obj); err != nil {
		return payload
	}

	maskSet := make(map[string]struct{}, len(t.cfg.MaskFields))
	for _, f := range t.cfg.MaskFields {
		maskSet[strings.TrimSpace(f)] = struct{}{}
	}

	for k, newKey := range t.cfg.RenameFields {
		if v, ok := obj[k]; ok {
			obj[newKey] = v
			delete(obj, k)
		}
	}

	for k := range maskSet {
		if _, ok := obj[k]; ok {
			obj[k] = "***"
		}
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return payload
	}
	return string(out)
}
