// Package header provides header (property) extraction and matching for
// Apache Pulsar messages.
//
// Pulsar messages carry an optional map of string properties (headers)
// analogous to HTTP headers. This package allows callers to:
//
//   - Filter messages whose properties contain a key and/or value matching
//     a regular expression via [Extractor.Match].
//   - Obtain a subset of properties whose keys satisfy the key pattern via
//     [Extractor.Extract].
//
// Example:
//
//	e, err := header.New("^x-", "")
//	if err != nil {
//		log.Fatal(err)
//	}
//	if e.Match(msg.Properties()) {
//		fmt.Println(e.Extract(msg.Properties()))
//	}
package header
