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

// SubscriptionOptionFunc -
type SubscriptionOptionFunc func(opts *SubscriberOptions) error

// Key represents the Ctx's key
type Key string
