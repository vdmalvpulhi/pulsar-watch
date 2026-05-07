// Package topicmap provides a concurrent, size-bounded map that tracks
// per-topic message statistics (seen, matched, last-seen timestamp) for
// use in the pulsar-watch monitoring pipeline.
//
// Usage:
//
//	tm, err := topicmap.New(100) // track up to 100 topics
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	tm.RecordSeen("persistent://public/default/orders")
//	tm.RecordMatched("persistent://public/default/orders")
//
//	for _, entry := range tm.Snapshot() {
//		fmt.Printf("%s seen=%d matched=%d\n", entry.Topic, entry.Seen, entry.Matched)
//	}
//
// When maxKeys is 0 the map grows without bound. When maxKeys > 0 new
// topics that would exceed the limit are silently dropped, preserving
// memory in high-cardinality environments.
package topicmap
