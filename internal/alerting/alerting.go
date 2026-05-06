// Package alerting provides threshold-based alerting for message statistics.
// It evaluates a set of rules against a stats snapshot and fires callbacks
// when conditions are met.
package alerting

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/user/pulsar-watch/internal/stats"
)

// Condition represents the type of comparison used in a rule.
type Condition string

const (
	ConditionGT  Condition = "gt"  // greater than
	ConditionGTE Condition = "gte" // greater than or equal
	ConditionLT  Condition = "lt"  // less than
	ConditionLTE Condition = "lte" // less than or equal
)

// Rule defines a single alerting rule evaluated against a stats snapshot.
type Rule struct {
	Name      string
	Field     string    // "seen", "matched", "exported", "dropped"
	Condition Condition
	Threshold int64
}

// Alert is produced when a rule's condition is satisfied.
type Alert struct {
	Rule      Rule
	Value     int64
	FiredAt   time.Time
}

// Handler is called when an alert fires.
type Handler func(Alert)

// Alerter evaluates rules against stats snapshots.
type Alerter struct {
	mu      sync.Mutex
	rules   []Rule
	handler Handler
	fired   map[string]bool // tracks already-fired rules to avoid repeat alerts
}

// New creates an Alerter with the given rules and alert handler.
func New(rules []Rule, handler Handler) (*Alerter, error) {
	if handler == nil {
		return nil, errors.New("alerting: handler must not be nil")
	}
	for _, r := range rules {
		if r.Name == "" {
			return nil, errors.New("alerting: rule name must not be empty")
		}
		switch r.Field {
		case "seen", "matched", "exported", "dropped":
		default:
			return nil, fmt.Errorf("alerting: unknown field %q in rule %q", r.Field, r.Name)
		}
		switch r.Condition {
		case ConditionGT, ConditionGTE, ConditionLT, ConditionLTE:
		default:
			return nil, fmt.Errorf("alerting: unknown condition %q in rule %q", r.Condition, r.Name)
		}
	}
	return &Alerter{
		rules:   rules,
		handler: handler,
		fired:   make(map[string]bool),
	}, nil
}

// Evaluate checks all rules against the provided snapshot, firing the handler
// for each rule whose condition is met and has not yet fired.
func (a *Alerter) Evaluate(snap stats.Snapshot) {
	a.mu.Lock()
	defer a.mu.Unlock()
	for _, rule := range a.rules {
		if a.fired[rule.Name] {
			continue
		}
		val := fieldValue(snap, rule.Field)
		if matches(val, rule.Condition, rule.Threshold) {
			a.fired[rule.Name] = true
			a.handler(Alert{Rule: rule, Value: val, FiredAt: time.Now()})
		}
	}
}

// Reset clears all fired state so rules can fire again.
func (a *Alerter) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.fired = make(map[string]bool)
}

func fieldValue(snap stats.Snapshot, field string) int64 {
	switch field {
	case "seen":
		return snap.Seen
	case "matched":
		return snap.Matched
	case "exported":
		return snap.Exported
	case "dropped":
		return snap.Dropped
	}
	return 0
}

func matches(val int64, cond Condition, threshold int64) bool {
	switch cond {
	case ConditionGT:
		return val > threshold
	case ConditionGTE:
		return val >= threshold
	case ConditionLT:
		return val < threshold
	case ConditionLTE:
		return val <= threshold
	}
	return false
}
