package exporter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Format represents the supported export formats.
type Format string

const (
	FormatJSON Format = "json"
	FormatText Format = "text"
)

// Message represents a Pulsar message to be exported.
type Message struct {
	Topic     string            `json:"topic"`
	Key       string            `json:"key"`
	Payload   string            `json:"payload"`
	Timestamp time.Time         `json:"timestamp"`
	Headers   map[string]string `json:"headers,omitempty"`
}

// Exporter writes messages to an output destination.
type Exporter struct {
	format Format
	writer io.Writer
}

// New creates a new Exporter. If outputPath is empty, stdout is used.
func New(format Format, outputPath string) (*Exporter, error) {
	if format != FormatJSON && format != FormatText {
		return nil, fmt.Errorf("unsupported export format: %q", format)
	}

	var w io.Writer = os.Stdout
	if outputPath != "" {
		f, err := os.Create(outputPath)
		if err != nil {
			return nil, fmt.Errorf("opening output file: %w", err)
		}
		w = f
	}

	return &Exporter{format: format, writer: w}, nil
}

// NewWithWriter creates an Exporter that writes to the provided writer.
func NewWithWriter(format Format, w io.Writer) (*Exporter, error) {
	if format != FormatJSON && format != FormatText {
		return nil, fmt.Errorf("unsupported export format: %q", format)
	}
	return &Exporter{format: format, writer: w}, nil
}

// Write encodes and writes a single message to the output destination.
func (e *Exporter) Write(msg Message) error {
	switch e.format {
	case FormatJSON:
		return e.writeJSON(msg)
	case FormatText:
		return e.writeText(msg)
	default:
		return fmt.Errorf("unsupported format: %q", e.format)
	}
}

func (e *Exporter) writeJSON(msg Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshalling message: %w", err)
	}
	_, err = fmt.Fprintf(e.writer, "%s\n", data)
	return err
}

func (e *Exporter) writeText(msg Message) error {
	_, err := fmt.Fprintf(e.writer, "[%s] topic=%s key=%s payload=%s\n",
		msg.Timestamp.Format(time.RFC3339),
		msg.Topic,
		msg.Key,
		msg.Payload,
	)
	return err
}
