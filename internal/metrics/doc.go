// Package metrics provides a human-readable table renderer for pulsar-watch
// statistics snapshots.
//
// It wraps a [stats.Snapshot] and formats key counters — messages seen,
// matched, exported, and errors — into an aligned tabular layout suitable
// for terminal output.
//
// Basic usage:
//
//	p := metrics.New()          // writes to stdout
//	p.Print(snap)               // render a single snapshot
//	p.PrintDelta(prev, curr, d) // render rates over an interval
//
// Use [NewWithWriter] to redirect output to any [io.Writer], which is
// particularly useful in tests.
package metrics
