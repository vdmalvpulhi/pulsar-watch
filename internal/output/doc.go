// Package output provides terminal output formatting utilities for the
// pulsar-watch CLI tool.
//
// It wraps colored, leveled log-style printing (Info, Warn, Error, Debug)
// and a structured message display function used when consuming or replaying
// Pulsar topic messages.
//
// # Printer
//
// The [Printer] type is the main entry point. Create one with [New], passing
// a verbose flag to control whether Debug output is emitted.
//
// Example usage:
//
//	p := output.New(verbose)
//	p.Info("connected to broker %s", brokerURL)
//	p.Warn("no messages received within timeout")
//	p.Message(topic, msg.Key(), msg.Payload(), msg.PublishTime())
//
// # Color support
//
// Color output is automatically disabled when writing to a non-TTY writer
// (e.g. a file or pipe) via the github.com/fatih/color package. It can also
// be disabled explicitly by calling color.NoColor = true before creating a
// Printer.
package output
