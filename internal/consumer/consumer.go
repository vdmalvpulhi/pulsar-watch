// Package consumer wraps the Apache Pulsar client to provide a simplified
// message-consumption interface used by the watcher, replay, and tail
// subsystems.
package consumer

import (
	"context"
	"fmt"
	"time"
)

// Message represents a single Pulsar message.
type Message interface {
	Key() string
	Payload() []byte
	PublishTime() time.Time
	Properties() map[string]string
}

// Options configures a Consumer.
type Options struct {
	BrokerURL    string
	Topic        string
	Subscription string
}

// Consumer wraps a Pulsar client consumer.
type Consumer struct {
	opts Options
}

// New creates a Consumer and validates options.
// A real implementation would dial the Pulsar broker here.
func New(opts Options) (*Consumer, error) {
	if opts.BrokerURL == "" {
		return nil, fmt.Errorf("consumer: broker URL is required")
	}
	if opts.Topic == "" {
		return nil, fmt.Errorf("consumer: topic is required")
	}
	if opts.Subscription == "" {
		opts.Subscription = "pulsar-watch"
	}
	return &Consumer{opts: opts}, nil
}

// Receive blocks until a message is available or ctx is cancelled.
func (c *Consumer) Receive(ctx context.Context) (Message, error) {
	// Real implementation delegates to the Pulsar client.
	<-ctx.Done()
	return nil, ctx.Err()
}

// Ack acknowledges a message so the broker advances the subscription cursor.
func (c *Consumer) Ack(_ Message) {}

// Close releases all resources held by the consumer.
func (c *Consumer) Close() {}

// BrokerURL returns the configured broker URL.
func (c *Consumer) BrokerURL() string { return c.opts.BrokerURL }

// Topic returns the configured topic.
func (c *Consumer) Topic() string { return c.opts.Topic }

// Subscription returns the subscription name in use.
func (c *Consumer) Subscription() string { return c.opts.Subscription }
