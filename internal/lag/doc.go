// Package lag provides consumer lag tracking for Apache Pulsar topics.
//
// A Tracker records the message backlog (number of unacknowledged messages)
// for each topic/subscription pair observed by pulsar-watch. Entries are
// updated in-place on successive calls to Record and can be retrieved
// individually via Get or in bulk via Snapshot.
//
// A Printer formats lag snapshots as a human-readable table suitable for
// terminal output, sorting entries by topic and subscription name.
//
// Example usage:
//
//	tr := lag.New()
//	_ = tr.Record("persistent://public/default/orders", "my-sub", 42)
//
//	p := lag.NewPrinter(os.Stdout)
//	p.Print(tr.Snapshot())
package lag
