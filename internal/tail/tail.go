// Package tail provides live-tail functionality for Pulsar topics,
// printing messages to an output writer as they arrive.
package tail

import (
	"context"
	"fmt"
	"time"

	"github.com/user/pulsar-watch/internal/consumer"
	"github.com/user/pulsar-watch/internal/filter"
	"github.com/user/pulsar-watch/internal/output"
	"github.com/user/pulsar-watch/internal/ratelimit"
	"github.com/user/pulsar-watch/internal/stats"
)

// Options configures a Tail session.
type Options struct {
	Consumer  *consumer.Consumer
	Filter    *filter.Filter
	Output    *output.Output
	Stats     *stats.Stats
	RateLimit *ratelimit.RateLimiter
	// MaxMessages stops tailing after N messages (0 = unlimited).
	MaxMessages int
	// ShowTimestamp prefixes each line with the message publish time.
	ShowTimestamp bool
}

// Tail streams messages from a Pulsar topic to the output writer.
type Tail struct {
	opts Options
}

// New creates a new Tail with the provided options.
func New(opts Options) (*Tail, error) {
	if opts.Consumer == nil {
		return nil, fmt.Errorf("tail: consumer is required")
	}
	if opts.Output == nil {
		return nil, fmt.Errorf("tail: output is required")
	}
	return &Tail{opts: opts}, nil
}

// Run starts tailing until ctx is cancelled or MaxMessages is reached.
func (t *Tail) Run(ctx context.Context) error {
	count := 0
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		msg, err := t.opts.Consumer.Receive(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			return fmt.Errorf("tail: receive error: %w", err)
		}

		if t.opts.Stats != nil {
			t.opts.Stats.RecordSeen()
		}

		if t.opts.Filter != nil && !t.opts.Filter.Match(msg) {
			t.opts.Consumer.Ack(msg)
			continue
		}

		if t.opts.Stats != nil {
			t.opts.Stats.RecordMatched()
		}

		if t.opts.RateLimit != nil {
			if err := t.opts.RateLimit.Wait(ctx); err != nil {
				return err
			}
		}

		t.print(msg)
		t.opts.Consumer.Ack(msg)

		count++
		if t.opts.MaxMessages > 0 && count >= t.opts.MaxMessages {
			return nil
		}
	}
}

func (t *Tail) print(msg consumer.Message) {
	if t.opts.ShowTimestamp {
		ts := msg.PublishTime().Format(time.RFC3339)
		t.opts.Output.Info("[%s] key=%s payload=%s", ts, msg.Key(), string(msg.Payload()))
		return
	}
	t.opts.Output.Info("key=%s payload=%s", msg.Key(), string(msg.Payload()))
}
