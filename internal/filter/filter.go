package filter

import (
	"regexp"
	"strings"
)

// Message represents a Pulsar message to be evaluated against filters.
type Message struct {
	Key     string
	Payload string
	Properties map[string]string
}

// Options holds the filtering criteria for messages.
type Options struct {
	KeyPattern     string
	PayloadPattern string
	PropertyKey    string
	PropertyValue  string
}

// Filter evaluates messages against configured criteria.
type Filter struct {
	keyRegex     *regexp.Regexp
	payloadRegex *regexp.Regexp
	opts         Options
}

// New creates a new Filter from the given Options.
// Returns an error if any regex pattern is invalid.
func New(opts Options) (*Filter, error) {
	f := &Filter{opts: opts}

	if opts.KeyPattern != "" {
		re, err := regexp.Compile(opts.KeyPattern)
		if err != nil {
			return nil, fmt.Errorf("invalid key pattern: %w", err)
		}
		f.keyRegex = re
	}

	if opts.PayloadPattern != "" {
		re, err := regexp.Compile(opts.PayloadPattern)
		if err != nil {
			return nil, fmt.Errorf("invalid payload pattern: %w", err)
		}
		f.payloadRegex = re
	}

	return f, nil
}

// Match returns true if the message satisfies all configured filter criteria.
func (f *Filter) Match(msg Message) bool {
	if f.keyRegex != nil && !f.keyRegex.MatchString(msg.Key) {
		return false
	}

	if f.payloadRegex != nil && !f.payloadRegex.MatchString(msg.Payload) {
		return false
	}

	if f.opts.PropertyKey != "" {
		val, ok := msg.Properties[f.opts.PropertyKey]
		if !ok {
			return false
		}
		if f.opts.PropertyValue != "" && !strings.EqualFold(val, f.opts.PropertyValue) {
			return false
		}
	}

	return true
}
