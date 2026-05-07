// Package stats provides lightweight, thread-safe tracking of message
// consumption metrics for a pulsar-watch session.
//
// It records three counters:
//
//   - TotalSeen:     every message received from the Pulsar topic.
//   - TotalMatched:  messages that passed the active filter rules.
//   - TotalExported: messages successfully written by the exporter.
//
// All counters are updated atomically, so RecordSeen, RecordMatched, and
// RecordExported are safe to call concurrently from multiple goroutines.
//
// Usage:
//
//	s := stats.New()
//
//	// in your consume loop:
//	s.RecordSeen()
//	if filter.Match(msg) {
//		s.RecordMatched()
//		exporter.Write(msg)
//		s.RecordExported()
//	}
//
//	// print a summary at the end:
//	fmt.Println(s.Snapshot())
//
// Snapshot returns a point-in-time copy of all counters as a [Snapshot] value,
// which can be logged, serialised, or compared without holding any lock.
package stats
