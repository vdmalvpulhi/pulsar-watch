// Package sampler provides probabilistic message sampling for pulsar-watch.
//
// A Sampler is configured with a rate in [0.0, 1.0] where 0.0 drops every
// message and 1.0 keeps every message. Values in between cause each message
// to be independently accepted with the given probability.
//
// Example usage:
//
//	s, err := sampler.New(0.1) // keep ~10 % of messages
//	if err != nil {
//		log.Fatal(err)
//	}
//	if s.Sample() {
//		// process message
//	}
//
// Stats() can be used to report how many messages were seen versus kept,
// which is surfaced in the metrics output.
package sampler
