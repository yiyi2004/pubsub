package main

import (
	pubsub "github.com/DemonDCC/pubsub/mq/nats"
)

func main() {
	b := pubsub.NewBroker()

	p := b.CreatePublisher()

	p.Publish(b, "hello", []byte("I and demon"))

	b.Close()
}
