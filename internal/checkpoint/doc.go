// Package checkpoint provides persistent tracking of the last successfully
// processed Pulsar message ID per topic. Checkpoints are stored as a JSON
// file on disk and allow pulsar-watch to resume consumption from where it
// left off after a restart, avoiding duplicate processing or missed messages.
//
// Usage:
//
//	cp, err := checkpoint.New("/var/lib/pulsar-watch/checkpoints.json")
//	if err != nil { ... }
//
//	// Save progress after processing a message.
//	cp.Set("persistent://public/default/events", msg.ID)
//
//	// Retrieve the last known position on startup.
//	if id, ok := cp.Get("persistent://public/default/events"); ok {
//		// seek consumer to id
//	}
package checkpoint
