// Package header provides utilities for extracting and matching
// Apache Pulsar message headers (properties).
package header

import (
	"errors"
	"fmt"
	"regexp"
)

// Extractor holds compiled patterns for header key/value matching.
type Extractor struct {
	keyPattern   *regexp.Regexp
	valuePattern *regexp.Regexp
}

// New creates an Extractor that filters message properties by key and/or
// value regular expressions. Either pattern may be empty to skip that check.
func New(keyPattern, valuePattern string) (*Extractor, error) {
	var kp, vp *regexp.Regexp
	var err error

	if keyPattern != "" {
		kp, err = regexp.Compile(keyPattern)
		if err != nil {
			return nil, fmt.Errorf("header: invalid key pattern: %w", err)
		}
	}

	if valuePattern != "" {
		vp, err = regexp.Compile(valuePattern)
		if err != nil {
			return nil, fmt.Errorf("header: invalid value pattern: %w", err)
		}
	}

	if kp == nil && vp == nil {
		return nil, errors.New("header: at least one of keyPattern or valuePattern must be set")
	}

	return &Extractor{keyPattern: kp, valuePattern: vp}, nil
}

// Match returns true when the supplied properties map satisfies the
// configured key and value patterns.
//
// If only a key pattern is set, at least one key must match.
// If only a value pattern is set, at least one value must match.
// If both are set, a single entry must satisfy both simultaneously.
func (e *Extractor) Match(props map[string]string) bool {
	if e.keyPattern != nil && e.valuePattern != nil {
		for k, v := range props {
			if e.keyPattern.MatchString(k) && e.valuePattern.MatchString(v) {
				return true
			}
		}
		return false
	}

	if e.keyPattern != nil {
		for k := range props {
			if e.keyPattern.MatchString(k) {
				return true
			}
		}
		return false
	}

	// valuePattern only
	for _, v := range props {
		if e.valuePattern.MatchString(v) {
			return true
		}
	}
	return false
}

// Extract returns a new map containing only the entries whose keys match
// the key pattern. If no key pattern is configured all entries are returned.
func (e *Extractor) Extract(props map[string]string) map[string]string {
	out := make(map[string]string, len(props))
	for k, v := range props {
		if e.keyPattern == nil || e.keyPattern.MatchString(k) {
			out[k] = v
		}
	}
	return out
}
