// Package output provides terminal output formatting utilities for the
// pulsar-watch CLI tool.
//
// It wraps colored, leveled log-style printing (Info, Warn, Error, Debug)
// and a structured message display function used when consuming or replaying
// Pulsar topic messages.
//
// Example usage:
//
//	p := output.New(verbose)
//	p.Info("connected to broker %s", brokerURL)
//	p.Message(topic, msg.Key(), msg.Payload(), msg.PublishTime())
//
// Color output is automatically disabled when writing to a non-TTY writer
// (e.g. a file or pipe) via the github.com/fatih/color package.
package output
