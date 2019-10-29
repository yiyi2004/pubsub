package pubsub

import (
	"context"
	"errors"
)

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

// TopicOption -
func TopicOption(topic string) PublisherOptionFunc {
	return func(opts *PublisherOptions) error {
		if topic == "" {
			return errors.New("[error]: invalid topic")
		}

		opts.Ctx = context.WithValue(opts.Ctx, Key("topic"), topic)

		return nil
	}
}
