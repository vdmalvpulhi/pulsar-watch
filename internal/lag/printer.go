package lag

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// Printer formats and writes lag snapshots to an io.Writer.
type Printer struct {
	w io.Writer
}

// NewPrinter creates a Printer that writes to w.
func NewPrinter(w io.Writer) *Printer {
	return &Printer{w: w}
}

// Print writes a formatted table of lag entries to the writer.
func (p *Printer) Print(entries []Entry) {
	if len(entries) == 0 {
		fmt.Fprintln(p.w, "no lag data available")
		return
	}

	sorted := make([]Entry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Topic != sorted[j].Topic {
			return sorted[i].Topic < sorted[j].Topic
		}
		return sorted[i].Subscription < sorted[j].Subscription
	})

	tw := tabwriter.NewWriter(p.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TOPIC\tSUBSCRIPTION\tBACKLOG\tLAST UPDATED")
	for _, e := range sorted {
		fmt.Fprintf(tw, "%s\t%s\t%d\t%s\n",
			e.Topic,
			e.Subscription,
			e.Backlog,
			e.LastUpdated.Format("15:04:05"),
		)
	}
	_ = tw.Flush()
}
