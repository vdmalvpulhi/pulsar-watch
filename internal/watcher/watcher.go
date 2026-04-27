package watcher

import (
	"context"
	"fmt"

	"github.com/user/pulsar-watch/internal/consumer"
	"github.com/user/pulsar-watch/internal/exporter"
	"github.com/user/pulsar-watch/internal/filter"
	"github.com/user/pulsar-watch/internal/output"
	"github.com/user/pulsar-watch/internal/stats"
)

// Watcher coordinates message consumption, filtering, and exporting.
type Watcher struct {
	consumer consumer.Consumer
	filter   *filter.Filter
	exporter *exporter.Exporter
	stats    *stats.Stats
	out      *output.Output
}

// Config holds the dependencies needed to create a Watcher.
type Config struct {
	Consumer consumer.Consumer
	Filter   *filter.Filter
	Exporter *exporter.Exporter
	Stats    *stats.Stats
	Output   *output.Output
}

// New creates a new Watcher from the provided Config.
// Returns an error if any required dependency is nil.
func New(cfg Config) (*Watcher, error) {
	if cfg.Consumer == nil {
		return nil, fmt.Errorf("watcher: consumer is required")
	}
	if cfg.Filter == nil {
		return nil, fmt.Errorf("watcher: filter is required")
	}
	if cfg.Stats == nil {
		return nil, fmt.Errorf("watcher: stats is required")
	}
	if cfg.Output == nil {
		return nil, fmt.Errorf("watcher: output is required")
	}
	return &Watcher{
		consumer: cfg.Consumer,
		filter:   cfg.Filter,
		exporter: cfg.Exporter,
		stats:    cfg.Stats,
		out:      cfg.Output,
	}, nil
}

// Run starts the watch loop, consuming messages until the context is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	w.out.Info("starting watcher")
	for {
		select {
		case <-ctx.Done():
			w.out.Info("watcher stopped")
			return nil
		default:
		}

		msg, err := w.consumer.Receive(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			w.out.Warn(fmt.Sprintf("receive error: %v", err))
			continue
		}

		w.stats.RecordSeen()
		w.out.Debug(fmt.Sprintf("received message key=%s", msg.Key()))

		if !w.filter.Match(msg) {
			w.consumer.Ack(msg)
			continue
		}

		w.stats.RecordMatched()
		w.out.Info(fmt.Sprintf("matched message key=%s", msg.Key()))

		if w.exporter != nil {
			if err := w.exporter.Write(msg); err != nil {
				w.out.Warn(fmt.Sprintf("export error: %v", err))
			} else {
				w.stats.RecordExported()
			}
		}

		w.consumer.Ack(msg)
	}
}
