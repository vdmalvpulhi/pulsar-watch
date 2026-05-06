// Package throttle implements a sliding-window rate limiter for
// Pulsar message throughput control.
//
// A Throttle tracks message arrival times within a configurable time
// window and signals when the observed rate exceeds a configured
// messages-per-second ceiling. Unlike ratelimit (which blocks), throttle
// is advisory — callers decide how to react (e.g. pause, drop, or back off).
//
// Example:
//
//	t, _ := throttle.New(100, time.Second)
//	for msg := range messages {
//		if t.Record(time.Now()) {
//			// rate exceeded — slow down
//			time.Sleep(50 * time.Millisecond)
//		}
//		process(msg)
//	}
package throttle
