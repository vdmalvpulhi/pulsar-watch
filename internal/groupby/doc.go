// Package groupby provides message aggregation by an arbitrary string key,
// such as a Pulsar message key, topic name, or extracted payload field.
//
// A Grouper tracks seen and matched counts per key up to a configurable
// maximum number of distinct keys. Once the cap is reached new keys are
// silently dropped so that memory usage remains bounded.
//
// Usage:
//
//	g, err := groupby.New(1000)
//	if err != nil { ... }
//
//	g.RecordSeen(msg.Key)
//	if filter.Match(msg) {
//		g.RecordMatched(msg.Key)
//	}
//
//	p := groupby.NewPrinter(os.Stdout)
//	p.Print(g.Snapshot())
package groupby
