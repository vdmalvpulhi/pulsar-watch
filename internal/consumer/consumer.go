package consumer

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

// Message wraps a Pulsar message with metadata.
type Message struct {
	ID        string
	Key       string
	Payload   []byte
	Topic     string
	PublishTime time.Time
	Properties  map[string]string
}

// Consumer reads messages from a Pulsar topic.
type Consumer struct {
	client   pulsar.Client
	consumer pulsar.Consumer
}

// Options holds configuration for creating a Consumer.
type Options struct {
	BrokerURL      string
	Topic          string
	Subscription   string
	InitialPosition pulsar.SubscriptionInitialPosition
}

// New creates a new Consumer connected to the given broker.
func New(opts Options) (*Consumer, error) {
	if opts.BrokerURL == "" {
		return nil, fmt.Errorf("broker URL must not be empty")
	}
	if opts.Topic == "" {
		return nil, fmt.Errorf("topic must not be empty")
	}
	if opts.Subscription == "" {
		opts.Subscription = "pulsar-watch"
	}

	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: opts.BrokerURL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create pulsar client: %w", err)
	}

	c, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:                       opts.Topic,
		SubscriptionName:            opts.Subscription,
		Type:                        pulsar.Exclusive,
		SubscriptionInitialPosition: opts.InitialPosition,
	})
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to subscribe to topic: %w", err)
	}

	return &Consumer{client: client, consumer: c}, nil
}

// Receive blocks until a message is available or the context is cancelled.
func (c *Consumer) Receive(ctx context.Context) (*Message, error) {
	msg, err := c.consumer.Receive(ctx)
	if err != nil {
		return nil, fmt.Errorf("receive error: %w", err)
	}
	c.consumer.Ack(msg)
	return &Message{
		ID:          msg.ID().String(),
		Key:         msg.Key(),
		Payload:     msg.Payload(),
		Topic:       msg.Topic(),
		PublishTime: msg.PublishTime(),
		Properties:  msg.Properties(),
	}, nil
}

// Close releases resources held by the consumer.
func (c *Consumer) Close() {
	c.consumer.Close()
	c.client.Close()
}
