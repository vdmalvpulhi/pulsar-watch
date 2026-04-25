package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration for pulsar-watch.
type Config struct {
	Pulsar  PulsarConfig  `yaml:"pulsar"`
	Watch   WatchConfig   `yaml:"watch"`
	Export  ExportConfig  `yaml:"export"`
}

// PulsarConfig contains Pulsar broker connection settings.
type PulsarConfig struct {
	BrokerURL        string        `yaml:"broker_url"`
	AuthToken        string        `yaml:"auth_token"`
	ConnectionTimeout time.Duration `yaml:"connection_timeout"`
	OperationTimeout  time.Duration `yaml:"operation_timeout"`
}

// WatchConfig contains topic watching and filtering settings.
type WatchConfig struct {
	Topic          string   `yaml:"topic"`
	SubscriptionName string `yaml:"subscription_name"`
	Filters        []string `yaml:"filters"`
	MaxMessages    int      `yaml:"max_messages"`
}

// ExportConfig contains message export settings.
type ExportConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Format   string `yaml:"format"`   // json, csv, raw
	FilePath string `yaml:"file_path"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Pulsar: PulsarConfig{
			BrokerURL:         "pulsar://localhost:6650",
			ConnectionTimeout: 10 * time.Second,
			OperationTimeout:  30 * time.Second,
		},
		Watch: WatchConfig{
			SubscriptionName: "pulsar-watch-sub",
			MaxMessages:      0, // 0 means unlimited
		},
		Export: ExportConfig{
			Enabled: false,
			Format:  "json",
		},
	}
}

// Load reads a YAML config file from the given path and returns a Config.
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

// Validate checks that required fields are set and values are acceptable.
func (c *Config) Validate() error {
	if c.Pulsar.BrokerURL == "" {
		return fmt.Errorf("pulsar.broker_url must not be empty")
	}
	validFormats := map[string]bool{"json": true, "csv": true, "raw": true}
	if c.Export.Enabled && !validFormats[c.Export.Format] {
		return fmt.Errorf("export.format %q is not valid; choose json, csv, or raw", c.Export.Format)
	}
	return nil
}
