// Package exporter provides functionality for writing Pulsar messages
// to various output formats and destinations.
//
// Supported formats:
//
//   - json: each message is serialised as a JSON object on its own line
//     (newline-delimited JSON / NDJSON)
//
//   - text: each message is written as a human-readable single line
//     containing the timestamp, topic, key and payload
//
// Usage:
//
//	e, err := exporter.New(exporter.FormatJSON, "/tmp/messages.ndjson")
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = e.Write(exporter.Message{
//		Topic:     "persistent://public/default/events",
//		Key:       "user-123",
//		Payload:   `{"action":"login"}`,
//		Timestamp: time.Now(),
//	})
//
// If no output path is provided to New, messages are written to stdout.
package exporter
