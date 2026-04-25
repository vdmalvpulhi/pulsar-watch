// Package filter provides message filtering capabilities for pulsar-watch.
//
// It supports matching Pulsar messages against configurable criteria including:
//   - Key regex patterns
//   - Payload regex patterns
//   - Message property key/value pairs (case-insensitive value comparison)
//
// Example usage:
//
//	opts := filter.Options{
//		KeyPattern:    "^order-",
//		PayloadPattern: `"status":"failed"`,
//		PropertyKey:   "region",
//		PropertyValue: "us-east",
//	}
//
//	f, err := filter.New(opts)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if f.Match(msg) {
//		// process matching message
//	}
package filter
