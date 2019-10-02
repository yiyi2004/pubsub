package main

import (
	pubsub "github.com/zhangce1999/pubsub/mq/nats"
)

func main() {
	b := pubsub.NewBroker()

	p := b.CreatePublisher()

	p.Publish(b, "hello", []byte("I am demon"))

	b.Close()
}
