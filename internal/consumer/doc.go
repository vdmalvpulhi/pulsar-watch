// Package consumer provides a thin wrapper around the Apache Pulsar
// client library to simplify consuming messages from a single topic.
//
// Usage:
//
//	c, err := consumer.New(consumer.Options{
//		BrokerURL:       "pulsar://localhost:6650",
//		Topic:           "persistent://public/default/my-topic",
//		Subscription:    "pulsar-watch",
//		InitialPosition: pulsar.SubscriptionPositionEarliest,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer c.Close()
//
//	for {
//		msg, err := c.Receive(ctx)
//		if err != nil {
//			break
//		}
//		fmt.Println(string(msg.Payload))
//	}
package consumer
