package output_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/pulsar-watch/internal/output"
)

func newBuf(verbose bool) (*output.Printer, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	return output.NewWithWriter(buf, verbose), buf
}

func TestInfo_ContainsMessage(t *testing.T) {
	p, buf := newBuf(false)
	p.Info("hello %s", "world")
	if !strings.Contains(buf.String(), "hello world") {
		t.Errorf("expected 'hello world' in output, got: %s", buf.String())
	}
}

func TestWarn_ContainsMessage(t *testing.T) {
	p, buf := newBuf(false)
	p.Warn("something off")
	if !strings.Contains(buf.String(), "something off") {
		t.Errorf("expected warn message in output, got: %s", buf.String())
	}
}

func TestError_ContainsMessage(t *testing.T) {
	p, buf := newBuf(false)
	p.Error("boom %d", 42)
	if !strings.Contains(buf.String(), "boom 42") {
		t.Errorf("expected error message in output, got: %s", buf.String())
	}
}

func TestDebug_SuppressedWhenNotVerbose(t *testing.T) {
	p, buf := newBuf(false)
	p.Debug("secret")
	if buf.Len() != 0 {
		t.Errorf("expected no output in non-verbose mode, got: %s", buf.String())
	}
}

func TestDebug_ShownWhenVerbose(t *testing.T) {
	p, buf := newBuf(true)
	p.Debug("visible debug")
	if !strings.Contains(buf.String(), "visible debug") {
		t.Errorf("expected debug message in verbose mode, got: %s", buf.String())
	}
}

func TestMessage_ContainsFields(t *testing.T) {
	p, buf := newBuf(false)
	now := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	p.Message("persistent://public/default/test", "key-1", []byte(`{"v":1}`), now)
	out := buf.String()
	for _, want := range []string{"persistent://public/default/test", "key-1", `{"v":1}`, "2024-06-01"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in message output, got: %s", want, out)
		}
	}
}

func TestSuccess_ContainsMessage(t *testing.T) {
	p, buf := newBuf(false)
	p.Success("done")
	if !strings.Contains(buf.String(), "done") {
		t.Errorf("expected success message in output, got: %s", buf.String())
	}
}
