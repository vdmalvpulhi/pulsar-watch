package timeline

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// Printer renders timeline entries to a writer in a human-readable table.
type Printer struct {
	w io.Writer
}

// NewPrinter returns a Printer that writes to w. If w is nil, os.Stdout is used.
func NewPrinter(w io.Writer) *Printer {
	if w == nil {
		w = os.Stdout
	}
	return &Printer{w: w}
}

// Print writes all timeline entries as a formatted table.
func (p *Printer) Print(t *Timeline) {
	entries := t.Entries()
	tw := tabwriter.NewWriter(p.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TIME\tTOPIC\tKEY\tPAYLOAD")
	fmt.Fprintln(tw, "----\t-----\t---\t-------")
	for _, e := range entries {
		payload := e.Payload
		if len(payload) > 60 {
			payload = payload[:57] + "..."
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
			e.RecordedAt.Format("15:04:05.000"),
			e.Topic,
			e.Key,
			payload,
		)
	}
	tw.Flush()
}
