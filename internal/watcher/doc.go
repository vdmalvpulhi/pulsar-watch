// Package watcher provides the core watch loop for pulsar-watch.
//
// A Watcher ties together a consumer, filter, exporter, stats tracker,
// and output logger into a single coordinated pipeline:
//
//  1. Messages are received from the Pulsar topic via the Consumer.
//  2. Each message is evaluated by the Filter; non-matching messages
//     are acknowledged and skipped.
//  3. Matching messages are optionally written by the Exporter.
//  4. Stats are updated at each stage (seen, matched, exported).
//
// Example usage:
//
//	w, err := watcher.New(watcher.Config{
//		Consumer: c,
//		Filter:   f,
//		Exporter: exp,
//		Stats:    s,
//		Output:   out,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	w.Run(ctx)
package watcher
