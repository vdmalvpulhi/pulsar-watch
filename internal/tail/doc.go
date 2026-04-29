// Package tail implements live-tail streaming for Apache Pulsar topics.
//
// It consumes messages in real time, applies optional filtering and rate
// limiting, records statistics, and writes matching messages to the
// configured output writer.
//
// Basic usage:
//
//	t, err := tail.New(tail.Options{
//		Consumer:      c,
//		Filter:        f,
//		Output:        out,
//		ShowTimestamp: true,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	if err := t.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
//		log.Fatal(err)
//	}
package tail
