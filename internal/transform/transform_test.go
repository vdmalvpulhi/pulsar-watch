package transform

import (
	"strings"
	"testing"
)

func TestNew_NegativeTruncate(t *testing.T) {
	_, err := New(Config{TruncateAt: -1})
	if err == nil {
		t.Fatal("expected error for negative TruncateAt")
	}
}

func TestNew_Valid(t *testing.T) {
	tr, err := New(Config{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil Transformer")
	}
}

func TestApply_NonJSON_Passthrough(t *testing.T) {
	tr, _ := New(Config{MaskFields: []string{"secret"}})
	input := "plain text payload"
	if got := tr.Apply(input); got != input {
		t.Errorf("expected %q, got %q", input, got)
	}
}

func TestApply_MaskField(t *testing.T) {
	tr, _ := New(Config{MaskFields: []string{"password"}})
	input := `{"user":"alice","password":"s3cr3t"}`
	out := tr.Apply(input)
	if strings.Contains(out, "s3cr3t") {
		t.Errorf("masked field value should not appear in output: %s", out)
	}
	if !strings.Contains(out, "***") {
		t.Errorf("expected mask placeholder in output: %s", out)
	}
	if !strings.Contains(out, "alice") {
		t.Errorf("non-masked field should remain: %s", out)
	}
}

func TestApply_RenameField(t *testing.T) {
	tr, _ := New(Config{RenameFields: map[string]string{"old_key": "new_key"}})
	input := `{"old_key":"value"}`
	out := tr.Apply(input)
	if strings.Contains(out, "old_key") {
		t.Errorf("original key should be removed: %s", out)
	}
	if !strings.Contains(out, "new_key") {
		t.Errorf("renamed key should appear: %s", out)
	}
}

func TestApply_Truncate(t *testing.T) {
	tr, _ := New(Config{TruncateAt: 10})
	input := "hello world this is a long payload"
	out := tr.Apply(input)
	if len([]rune(out)) > 12 { // 10 chars + ellipsis rune
		t.Errorf("output too long: %q", out)
	}
	if !strings.HasSuffix(out, "…") {
		t.Errorf("expected ellipsis suffix: %q", out)
	}
}

func TestApply_TruncateZeroMeansUnlimited(t *testing.T) {
	tr, _ := New(Config{TruncateAt: 0})
	input := strings.Repeat("a", 500)
	if got := tr.Apply(input); got != input {
		t.Errorf("expected unchanged payload when TruncateAt=0")
	}
}

func TestApply_MaskAndRenameAndTruncate(t *testing.T) {
	tr, _ := New(Config{
		MaskFields:   []string{"token"},
		RenameFields: map[string]string{"id": "user_id"},
		TruncateAt:   40,
	})
	input := `{"id":"123","token":"abc","name":"bob"}`
	out := tr.Apply(input)
	if strings.Contains(out, "abc") {
		t.Errorf("token value should be masked")
	}
	if strings.Contains(out, `"id"`) {
		t.Errorf("old key 'id' should be renamed")
	}
	if len(out) > 42 {
		t.Errorf("output should be truncated, got len=%d", len(out))
	}
}
