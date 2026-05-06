package groupby

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// Printer renders group snapshots to an io.Writer.
type Printer struct {
	w io.Writer
}

// NewPrinter creates a Printer that writes to w.
func NewPrinter(w io.Writer) *Printer {
	return &Printer{w: w}
}

// Print writes a formatted table of groups sorted by Seen count descending.
func (p *Printer) Print(groups []Group) {
	if len(groups) == 0 {
		fmt.Fprintln(p.w, "no groups recorded")
		return
	}

	sorted := make([]Group, len(groups))
	copy(sorted, groups)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Seen > sorted[j].Seen
	})

	tw := tabwriter.NewWriter(p.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "KEY\tSEEN\tMATCHED\tLAST SEEN")
	fmt.Fprintln(tw, "---\t----\t-------\t---------")
	for _, g := range sorted {
		fmt.Fprintf(tw, "%s\t%d\t%d\t%s\n",
			g.Key,
			g.Seen,
			g.Matched,
			g.LastSeen.Format("15:04:05"),
		)
	}
	_ = tw.Flush()
}
