package schema_test

import (
	"strings"
	"testing"

	"github.com/user/pulsar-watch/internal/schema"
)

func TestDetect_EmptyPayload(t *testing.T) {
	if got := schema.Detect(nil); got != schema.TypeUnknown {
		t.Fatalf("expected Unknown, got %s", got)
	}
}

func TestDetect_WhitespaceOnly(t *testing.T) {
	if got := schema.Detect([]byte("   \n")); got != schema.TypeUnknown {
		t.Fatalf("expected Unknown, got %s", got)
	}
}

func TestDetect_ValidJSONObject(t *testing.T) {
	payload := []byte(`{"key":"value"}`)
	if got := schema.Detect(payload); got != schema.TypeJSON {
		t.Fatalf("expected JSON, got %s", got)
	}
}

func TestDetect_ValidJSONArray(t *testing.T) {
	payload := []byte(`[1,2,3]`)
	if got := schema.Detect(payload); got != schema.TypeJSON {
		t.Fatalf("expected JSON, got %s", got)
	}
}

func TestDetect_InvalidJSON(t *testing.T) {
	payload := []byte(`{not valid json}`)
	if got := schema.Detect(payload); got != schema.TypeText {
		t.Fatalf("expected Text, got %s", got)
	}
}

func TestDetect_XMLPayload(t *testing.T) {
	payload := []byte(`<root><item>value</item></root>`)
	if got := schema.Detect(payload); got != schema.TypeXML {
		t.Fatalf("expected XML, got %s", got)
	}
}

func TestDetect_PlainText(t *testing.T) {
	payload := []byte(`hello world`)
	if got := schema.Detect(payload); got != schema.TypeText {
		t.Fatalf("expected Text, got %s", got)
	}
}

func TestPretty_JSONIndented(t *testing.T) {
	payload := []byte(`{"a":1,"b":2}`)
	out, err := schema.Pretty(payload, schema.TypeJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "\n") {
		t.Fatal("expected indented JSON output")
	}
}

func TestPretty_InvalidJSONReturnsError(t *testing.T) {
	payload := []byte(`{bad}`)
	_, err := schema.Pretty(payload, schema.TypeJSON)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestPretty_TextPassthrough(t *testing.T) {
	payload := []byte(`hello`)
	out, err := schema.Pretty(payload, schema.TypeText)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "hello" {
		t.Fatalf("expected passthrough, got %q", out)
	}
}

func TestPretty_UnknownPassthrough(t *testing.T) {
	payload := []byte(`raw bytes`)
	out, err := schema.Pretty(payload, schema.TypeUnknown)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "raw bytes" {
		t.Fatalf("expected passthrough, got %q", out)
	}
}
