// Package offset provides a file-backed tracker for Apache Pulsar message
// sequence IDs, enabling pulsar-watch to resume consumption from the last
// observed position on a per-topic basis.
//
// Usage:
//
//	tr, err := offset.New("/var/lib/pulsar-watch/offsets.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Record the latest offset after processing a message.
//	_ = tr.Set("persistent://public/default/events", msg.ID().SequenceID())
//
//	// Persist to disk before exit.
//	_ = tr.Save()
//
//	// On next startup, resume from the saved position.
//	if seq, ok := tr.Get("persistent://public/default/events"); ok {
//		fmt.Println("resuming from", seq)
//	}
package offset
