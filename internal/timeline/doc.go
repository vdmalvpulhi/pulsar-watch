// Package timeline provides a bounded, thread-safe chronological log of
// Pulsar message events observed during a watch session.
//
// A Timeline records entries up to a configurable maximum size, evicting the
// oldest entry when full. Entries can be retrieved as a snapshot or printed
// to any io.Writer via the Printer helper.
//
// Example usage:
//
//	tl, err := timeline.New(500)
//	if err != nil {
//		log.Fatal(err)
//	}
//	tl.Add("persistent://public/default/events", "user-1", `{"action":"login"}`)
//	timeline.NewPrinter(os.Stdout).Print(tl)
package timeline
