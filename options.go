package pubsub

import "context"

// Options -
type Options struct {
	Ctx context.Context
}

// Option is a function to register the option for the broker, publisher and subscription
type Option func(opts *Options) error
