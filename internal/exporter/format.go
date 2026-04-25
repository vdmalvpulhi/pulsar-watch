package exporter

import (
	"fmt"
	"strings"
)

// ParseFormat converts a raw string into a validated Format value.
// It is case-insensitive and trims surrounding whitespace.
func ParseFormat(raw string) (Format, error) {
	normalised := strings.TrimSpace(strings.ToLower(raw))
	switch Format(normalised) {
	case FormatJSON:
		return FormatJSON, nil
	case FormatText:
		return FormatText, nil
	default:
		return "", fmt.Errorf("unknown export format %q: must be one of [json, text]", raw)
	}
}

// SupportedFormats returns a slice of all supported Format values.
func SupportedFormats() []Format {
	return []Format{FormatJSON, FormatText}
}

// String implements the fmt.Stringer interface.
func (f Format) String() string {
	return string(f)
}
