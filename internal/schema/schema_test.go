package schema_test

import (
	"strings"
	"testing"

	"github.com/user/pulsar-watch/internal/schema"
)

func TestDetect_EmptyPayload(t *testing.T) {
	if got := schema.Detect([]byte{}); got != schema.FormatText {
		t.Fatalf("expected FormatText for empty payload, got %v", got)
	}
}

func TestDetect_WhitespaceOnly(t *testing.T) {
	if got := schema.Detect([]byte("   \n")); got != schema.FormatText {
		t.Fatalf("expected FormatText for whitespace, got %v", got)
	}
}

func TestDetect_ValidJSONObject(t *testing.T) {
	if got := schema.Detect([]byte(`{"key":"value"}`)); got != schema.FormatJSON {
		t.Fatalf("expected FormatJSON, got %v", got)
	}
}

func TestDetect_ValidJSONArray(t *testing.T) {
	if got := schema.Detect([]byte(`[1,2,3]`)); got != schema.FormatJSON {
		t.Fatalf("expected FormatJSON for array, got %v", got)
	}
}

func TestDetect_InvalidJSON(t *testing.T) {
	if got := schema.Detect([]byte(`{bad json`)); got != schema.FormatText {
		t.Fatalf("expected FormatText for invalid json, got %v", got)
	}
}

func TestDetect_PlainText(t *testing.T) {
	if got := schema.Detect([]byte("hello world")); got != schema.FormatText {
		t.Fatalf("expected FormatText for plain text, got %v", got)
	}
}

func TestFormat_String(t *testing.T) {
	cases := []struct {
		f    schema.Format
		want string
	}{
		{schema.FormatJSON, "json"},
		{schema.FormatText, "text"},
		{schema.FormatUnknown, "unknown"},
	}
	for _, tc := range cases {
		if got := tc.f.String(); got != tc.want {
			t.Errorf("Format(%d).String() = %q, want %q", tc.f, got, tc.want)
		}
	}
}

func TestPretty_JSONIsIndented(t *testing.T) {
	payload := []byte(`{"a":1,"b":"hello"}`)
	out, err := schema.Pretty(payload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "\n") {
		t.Errorf("expected indented JSON to contain newlines, got: %s", out)
	}
	if !strings.Contains(out, "  ") {
		t.Errorf("expected indented JSON to contain spaces, got: %s", out)
	}
}

func TestPretty_TextPassthrough(t *testing.T) {
	payload := []byte("simple message")
	out, err := schema.Pretty(payload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "simple message" {
		t.Errorf("expected passthrough, got: %q", out)
	}
}
