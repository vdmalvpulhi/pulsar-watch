// Package ratelimit implements a token-bucket rate limiter used to throttle
// the rate at which pulsar-watch processes or replays topic messages.
//
// Usage:
//
//	limiter, err := ratelimit.New(100) // 100 messages per second
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer limiter.Stop()
//
//	for _, msg := range messages {
//		if err := limiter.Wait(ctx); err != nil {
//			break
//		}
//		process(msg)
//	}
//
// A rate of 0 disables throttling entirely.
package ratelimit
