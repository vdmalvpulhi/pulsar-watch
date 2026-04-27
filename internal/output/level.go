package output

import "strings"

// Level represents the verbosity level for output filtering.
type Level int

const (
	// LevelInfo shows info, success, warn, and error messages.
	LevelInfo Level = iota
	// LevelDebug shows all messages including debug.
	LevelDebug
)

// ParseLevel converts a string to a Level.
// It returns LevelInfo for unrecognized values.
func ParseLevel(s string) Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return LevelDebug
	default:
		return LevelInfo
	}
}

// String returns the string representation of a Level.
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	default:
		return "info"
	}
}

// IsDebug returns true when the level includes debug output.
func (l Level) IsDebug() bool {
	return l >= LevelDebug
}
