package transform

import (
	"strings"
	"testing"
)

func stages(entries ...struct {
	Name string
	Cfg  Config
}) []struct {
	Name string
	Cfg  Config
} {
	return entries
}

func entry(name string, cfg Config) struct {
	Name string
	Cfg  Config
} {
	return struct {
		Name string
		Cfg  Config
	}{Name: name, Cfg: cfg}
}

func TestNewPipeline_EmptyName(t *testing.T) {
	_, err := NewPipeline(stages(entry("", Config{})))
	if err == nil {
		t.Fatal("expected error for empty stage name")
	}
}

func TestNewPipeline_InvalidStage(t *testing.T) {
	_, err := NewPipeline(stages(entry("bad", Config{TruncateAt: -5})))
	if err == nil {
		t.Fatal("expected error for invalid stage config")
	}
}

func TestNewPipeline_Empty(t *testing.T) {
	p, err := NewPipeline(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Stages()) != 0 {
		t.Fatal("expected no stages")
	}
}

func TestPipeline_Apply_SingleStage(t *testing.T) {
	p, _ := NewPipeline(stages(entry("mask", Config{MaskFields: []string{"secret"}})))
	input := `{"secret":"hidden","pub":"visible"}`
	out := p.Apply(input)
	if strings.Contains(out, "hidden") {
		t.Errorf("secret value should be masked: %s", out)
	}
}

func TestPipeline_Apply_MultiStage(t *testing.T) {
	p, _ := NewPipeline(stages(
		entry("rename", Config{RenameFields: map[string]string{"k": "key"}}),
		entry("truncate", Config{TruncateAt: 20}),
	))
	input := `{"k":"` + strings.Repeat("x", 50) + `"}`
	out := p.Apply(input)
	if len(out) > 22 {
		t.Errorf("output should be truncated, got len=%d: %s", len(out), out)
	}
}

func TestPipeline_Stages_ReturnsCopy(t *testing.T) {
	p, _ := NewPipeline(stages(entry("s1", Config{})))
	ss := p.Stages()
	ss[0].Name = "mutated"
	if p.Stages()[0].Name == "mutated" {
		t.Error("Stages() should return a copy, not a reference")
	}
}
