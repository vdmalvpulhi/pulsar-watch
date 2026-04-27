// Package cursor provides persistent read-position tracking for pulsar-watch.
//
// A Store maps topic names to their last successfully processed message
// position, allowing pulsar-watch to resume consumption after a restart
// without re-processing already-seen messages.
//
// Positions are serialised as JSON to a file on disk. The zero position
// (empty MessageID) indicates that no messages have been processed yet
// for that topic.
//
// Example usage:
//
//	store, err := cursor.New("/var/lib/pulsar-watch/cursors.json")
//	if err != nil { ... }
//
//	if pos, ok := store.Get("persistent://public/default/events"); ok {
//		fmt.Println("resuming from", pos.MessageID)
//	}
//
//	err = store.Set(cursor.Position{
//		Topic:       "persistent://public/default/events",
//		MessageID:   "1234:0",
//		PublishTime: time.Now(),
//	})
package cursor
