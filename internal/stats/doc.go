// Package stats provides lightweight, thread-safe tracking of message
// consumption metrics for a pulsar-watch session.
//
// It records three counters:
//
//   - TotalSeen:     every message received from the Pulsar topic.
//   - TotalMatched:  messages that passed the active filter rules.
//   - TotalExported: messages successfully written by the exporter.
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
package stats
