// Package schema provides utilities for detecting and pretty-printing
// message payload formats such as JSON, XML, and plain text.
//
// Detection is performed by inspecting the raw bytes of a message payload
// and attempting to parse it as a known format. If no known format is
// detected, the payload is treated as plain text.
//
// Usage:
//
//	format := schema.Detect(payload)
//	pretty, err := schema.Pretty(payload, format)
package schema
