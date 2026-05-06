package alerting_test

import (
	"sync"
	"testing"
	"time"

	"github.com/user/pulsar-watch/internal/alerting"
	"github.com/user/pulsar-watch/internal/stats"
)

func nopHandler(_ alerting.Alert) {}

func snap(seen, matched, exported, dropped int64) stats.Snapshot {
	return stats.Snapshot{
		Seen:          seen,
		Matched:       matched,
		Exported:      exported,
		Dropped:       dropped,
		LastMessageAt: time.Time{},
	}
}

func TestNew_NilHandler(t *testing.T) {
	_, err := alerting.New(nil, nil)
	if err == nil {
		t.Fatal("expected error for nil handler")
	}
}

func TestNew_EmptyRuleName(t *testing.T) {
	rules := []alerting.Rule{{Name: "", Field: "seen", Condition: alerting.ConditionGT, Threshold: 10}}
	_, err := alerting.New(rules, nopHandler)
	if err == nil {
		t.Fatal("expected error for empty rule name")
	}
}

func TestNew_UnknownField(t *testing.T) {
	rules := []alerting.Rule{{Name: "r", Field: "unknown", Condition: alerting.ConditionGT, Threshold: 1}}
	_, err := alerting.New(rules, nopHandler)
	if err == nil {
		t.Fatal("expected error for unknown field")
	}
}

func TestNew_UnknownCondition(t *testing.T) {
	rules := []alerting.Rule{{Name: "r", Field: "seen", Condition: "eq", Threshold: 1}}
	_, err := alerting.New(rules, nopHandler)
	if err == nil {
		t.Fatal("expected error for unknown condition")
	}
}

func TestEvaluate_FiresWhenConditionMet(t *testing.T) {
	var mu sync.Mutex
	var fired []alerting.Alert
	rules := []alerting.Rule{
		{Name: "high-seen", Field: "seen", Condition: alerting.ConditionGTE, Threshold: 100},
	}
	a, err := alerting.New(rules, func(al alerting.Alert) {
		mu.Lock()
		fired = append(fired, al)
		mu.Unlock()
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	a.Evaluate(snap(100, 0, 0, 0))
	mu.Lock()
	defer mu.Unlock()
	if len(fired) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(fired))
	}
	if fired[0].Rule.Name != "high-seen" {
		t.Errorf("unexpected rule name: %s", fired[0].Rule.Name)
	}
}

func TestEvaluate_DoesNotFireTwice(t *testing.T) {
	count := 0
	rules := []alerting.Rule{
		{Name: "r", Field: "matched", Condition: alerting.ConditionGT, Threshold: 5},
	}
	a, _ := alerting.New(rules, func(_ alerting.Alert) { count++ })
	a.Evaluate(snap(0, 10, 0, 0))
	a.Evaluate(snap(0, 20, 0, 0))
	if count != 1 {
		t.Errorf("expected rule to fire once, fired %d times", count)
	}
}

func TestEvaluate_DoesNotFireWhenConditionNotMet(t *testing.T) {
	count := 0
	rules := []alerting.Rule{
		{Name: "r", Field: "exported", Condition: alerting.ConditionGT, Threshold: 50},
	}
	a, _ := alerting.New(rules, func(_ alerting.Alert) { count++ })
	a.Evaluate(snap(0, 0, 10, 0))
	if count != 0 {
		t.Errorf("expected no alerts, got %d", count)
	}
}

func TestReset_AllowsRuleToFireAgain(t *testing.T) {
	count := 0
	rules := []alerting.Rule{
		{Name: "r", Field: "dropped", Condition: alerting.ConditionGTE, Threshold: 1},
	}
	a, _ := alerting.New(rules, func(_ alerting.Alert) { count++ })
	a.Evaluate(snap(0, 0, 0, 1))
	a.Reset()
	a.Evaluate(snap(0, 0, 0, 2))
	if count != 2 {
		t.Errorf("expected 2 alerts after reset, got %d", count)
	}
}
