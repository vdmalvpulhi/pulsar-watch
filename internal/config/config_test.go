package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Pulsar.BrokerURL != "pulsar://localhost:6650" {
		t.Errorf("expected default broker URL, got %q", cfg.Pulsar.BrokerURL)
	}
	if cfg.Pulsar.ConnectionTimeout != 10*time.Second {
		t.Errorf("expected 10s connection timeout, got %v", cfg.Pulsar.ConnectionTimeout)
	}
	if cfg.Watch.SubscriptionName != "pulsar-watch-sub" {
		t.Errorf("expected default subscription name, got %q", cfg.Watch.SubscriptionName)
	}
	if cfg.Export.Format != "json" {
		t.Errorf("expected default export format json, got %q", cfg.Export.Format)
	}
}

func TestLoad_FileNotExist(t *testing.T) {
	cfg, err := Load("/nonexistent/path/config.yaml")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if cfg == nil {
		t.Fatal("expected default config, got nil")
	}
}

func TestLoad_ValidYAML(t *testing.T) {
	content := `
pulsar:
  broker_url: "pulsar://broker:6650"
  auth_token: "mytoken"
watch:
  topic: "persistent://public/default/events"
  subscription_name: "my-sub"
  max_messages: 100
export:
  enabled: true
  format: "csv"
  file_path: "/tmp/out.csv"
`
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Pulsar.BrokerURL != "pulsar://broker:6650" {
		t.Errorf("broker URL mismatch: %q", cfg.Pulsar.BrokerURL)
	}
	if cfg.Watch.MaxMessages != 100 {
		t.Errorf("max_messages mismatch: %d", cfg.Watch.MaxMessages)
	}
	if cfg.Export.Format != "csv" {
		t.Errorf("export format mismatch: %q", cfg.Export.Format)
	}
}

func TestValidate_InvalidExportFormat(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Export.Enabled = true
	cfg.Export.Format = "xml"

	if err := cfg.Validate(); err == nil {
		t.Error("expected validation error for invalid export format, got nil")
	}
}

func TestValidate_EmptyBrokerURL(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Pulsar.BrokerURL = ""

	if err := cfg.Validate(); err == nil {
		t.Error("expected validation error for empty broker URL, got nil")
	}
}
