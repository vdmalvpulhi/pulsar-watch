// Package output provides terminal output formatting utilities
// for the pulsar-watch CLI, including colored and structured display.
package output

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fatih/color"
)

// Printer handles formatted output to a writer.
type Printer struct {
	w       io.Writer
	verbose bool
}

// New creates a Printer writing to stdout.
func New(verbose bool) *Printer {
	return &Printer{w: os.Stdout, verbose: verbose}
}

// NewWithWriter creates a Printer writing to the given writer.
func NewWithWriter(w io.Writer, verbose bool) *Printer {
	return &Printer{w: w, verbose: verbose}
}

// Info prints an informational message in cyan.
func (p *Printer) Info(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintln(p.w, color.CyanString("[INFO] ")+msg)
}

// Success prints a success message in green.
func (p *Printer) Success(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintln(p.w, color.GreenString("[OK]   ")+msg)
}

// Warn prints a warning message in yellow.
func (p *Printer) Warn(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintln(p.w, color.YellowString("[WARN] ")+msg)
}

// Error prints an error message in red.
func (p *Printer) Error(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintln(p.w, color.RedString("[ERR]  ")+msg)
}

// Debug prints a debug message only when verbose mode is enabled.
func (p *Printer) Debug(format string, args ...any) {
	if !p.verbose {
		return
	}
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintln(p.w, color.HiBlackString("[DBG]  ")+msg)
}

// Message prints a received Pulsar message in a human-readable format.
func (p *Printer) Message(topic, key string, payload []byte, publishTime time.Time) {
	header := color.MagentaString("[MSG]  ")
	fmt.Fprintf(p.w, "%stopic=%s key=%s time=%s\n",
		header,
		topic,
		key,
		publishTime.Format(time.RFC3339),
	)
	fmt.Fprintf(p.w, "       %s\n", string(payload))
}
