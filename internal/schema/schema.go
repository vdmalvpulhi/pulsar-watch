// Package schema provides message schema detection and pretty-printing
// for common payload formats such as JSON and plain text.
package schema

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Format represents a detected payload format.
type Format int

const (
	// FormatUnknown indicates the payload could not be classified.
	FormatUnknown Format = iota
	// FormatJSON indicates the payload is valid JSON.
	FormatJSON
	// FormatText indicates the payload is plain text.
	FormatText
)

// String returns a human-readable name for the format.
func (f Format) String() string {
	switch f {
	case FormatJSON:
		return "json"
	case FormatText:
		return "text"
	default:
		return "unknown"
	}
}

// Detect inspects the payload bytes and returns the best-guess Format.
func Detect(payload []byte) Format {
	trimmed := bytes.TrimSpace(payload)
	if len(trimmed) == 0 {
		return FormatText
	}
	if (trimmed[0] == '{' || trimmed[0] == '[') && json.Valid(trimmed) {
		return FormatJSON
	}
	return FormatText
}

// Pretty returns a formatted string representation of the payload.
// JSON payloads are indented; all others are returned as-is.
func Pretty(payload []byte) (string, error) {
	switch Detect(payload) {
	case FormatJSON:
		var buf bytes.Buffer
		if err := json.Indent(&buf, bytes.TrimSpace(payload), "", "  "); err != nil {
			return "", fmt.Errorf("schema: indent json: %w", err)
		}
		return buf.String(), nil
	default:
		return string(payload), nil
	}
}
