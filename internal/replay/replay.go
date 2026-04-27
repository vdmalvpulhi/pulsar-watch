// Package replay provides functionality to replay messages from a Pulsar topic
// with optional rate limiting and filtering support.
package replay

import (
	"context"
	"fmt"
	"time"

	"github.com/user/pulsar-watch/internal/consumer"
	"github.com/user/pulsar-watch/internal/exporter"
	"github.com/user/pulsar-watch/internal/filter"
)

// Options configures the replay behavior.
type Options struct {
	// RateLimit defines the maximum number of messages per second (0 = unlimited).
	RateLimit int
	// MaxMessages limits the total number of messages replayed (0 = unlimited).
	MaxMessages int
	// DryRun prints messages without exporting them.
	DryRun bool
}

// Replayer replays messages from a Pulsar topic.
type Replayer struct {
	consumer consumer.Consumer
	filter   *filter.Filter
	exporter *exporter.Exporter
	opts     Options
}

// New creates a new Replayer instance.
func New(c consumer.Consumer, f *filter.Filter, e *exporter.Exporter, opts Options) *Replayer {
	return &Replayer{
		consumer: c,
		filter:   f,
		exporter: e,
		opts:     opts,
	}
}

// Run starts replaying messages until the context is cancelled or MaxMessages is reached.
func (r *Replayer) Run(ctx context.Context) (int, error) {
	var (
		count   int
		ticker  *time.Ticker
	)

	if r.opts.RateLimit > 0 {
		interval := time.Second / time.Duration(r.opts.RateLimit)
		ticker = time.NewTicker(interval)
		defer ticker.Stop()
	}

	for {
		select {
		case <-ctx.Done():
			return count, nil
		default:
		}

		if r.opts.MaxMessages > 0 && count >= r.opts.MaxMessages {
			return count, nil
		}

		if ticker != nil {
			select {
			case <-ticker.C:
			case <-ctx.Done():
				return count, nil
			}
		}

		msg, err := r.consumer.Receive(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return count, nil
			}
			return count, fmt.Errorf("replay: receive message: %w", err)
		}

		if r.filter != nil && !r.filter.Match(msg) {
			r.consumer.Ack(msg)
			continue
		}

		if !r.opts.DryRun {
			if err := r.exporter.Write(msg); err != nil {
				return count, fmt.Errorf("replay: export message: %w", err)
			}
		}

		r.consumer.Ack(msg)
		count++
	}
}
