package pubsub

import "context"

// PublisherOptions -
type PublisherOptions struct {
	Ctx context.Context
}

// SubscriberOptions -
type SubscriberOptions struct {
	Ctx context.Context
}

// PublisherOptionFunc -
type PublisherOptionFunc func(opts *PublisherOptions) error

// SubscriberOptionFunc -
type SubscriberOptionFunc func(opts *SubscriberOptions) error
