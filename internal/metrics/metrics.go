package metrics

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	"github.com/user/pulsar-watch/internal/stats"
)

// Printer renders stats snapshots to a writer in a human-readable table.
type Printer struct {
	w io.Writer
}

// New returns a Printer that writes to stdout.
func New() *Printer {
	return &Printer{w: os.Stdout}
}

// NewWithWriter returns a Printer that writes to the given writer.
func NewWithWriter(w io.Writer) *Printer {
	return &Printer{w: w}
}

// Print renders a snapshot as a formatted metrics table.
func (p *Printer) Print(snap stats.Snapshot) error {
	tw := tabwriter.NewWriter(p.w, 0, 0, 2, ' ', 0)

	fmt.Fprintln(tw, "METRIC\tVALUE")
	fmt.Fprintln(tw, "------\t-----")
	fmt.Fprintf(tw, "Seen\t%d\n", snap.Seen)
	fmt.Fprintf(tw, "Matched\t%d\n", snap.Matched)
	fmt.Fprintf(tw, "Exported\t%d\n", snap.Exported)
	fmt.Fprintf(tw, "Errors\t%d\n", snap.Errors)

	if !snap.LastMessageAt.IsZero() {
		fmt.Fprintf(tw, "Last Message\t%s\n", snap.LastMessageAt.Format(time.RFC3339))
	} else {
		fmt.Fprintf(tw, "Last Message\t%s\n", "—")
	}

	return tw.Flush()
}

// PrintDelta renders the difference between two snapshots.
func (p *Printer) PrintDelta(prev, curr stats.Snapshot, interval time.Duration) error {
	tw := tabwriter.NewWriter(p.w, 0, 0, 2, ' ', 0)

	secs := interval.Seconds()
	rate := func(delta int64) float64 {
		if secs == 0 {
			return 0
		}
		return float64(delta) / secs
	}

	fmt.Fprintln(tw, "METRIC\tTOTAL\tRATE/s")
	fmt.Fprintln(tw, "------\t-----\t------")
	fmt.Fprintf(tw, "Seen\t%d\t%.2f\n", curr.Seen, rate(curr.Seen-prev.Seen))
	fmt.Fprintf(tw, "Matched\t%d\t%.2f\n", curr.Matched, rate(curr.Matched-prev.Matched))
	fmt.Fprintf(tw, "Exported\t%d\t%.2f\n", curr.Exported, rate(curr.Exported-prev.Exported))
	fmt.Fprintf(tw, "Errors\t%d\t%.2f\n", curr.Errors, rate(curr.Errors-prev.Errors))

	return tw.Flush()
}
