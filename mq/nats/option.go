package nats

import (
	"context"
	"log"

	"github.com/nats-io/nats.go"
)

// BrokerOptions -
type BrokerOptions struct {
	CustomedOptFunc []nats.Option
	// DefOpts represents a default nats Options
	UseDefault bool
	DefOpts    *nats.Options
	Ctx        context.Context
}

// BrokerOptionFunc -
type BrokerOptionFunc func(opts *BrokerOptions) error

// Connect -
func (bo *BrokerOptions) Connect() (*nats.Conn, error) {
	if bo.UseDefault {
		log.Println("Default Broker Options")
		return bo.DefOpts.Connect()
	}

	return nats.Connect(NATSURL, bo.CustomedOptFunc...)
}

// RegisterOptions -
func (bo *BrokerOptions) RegisterOptions(opts ...nats.Option) {
	if len(opts) == 0 {
		bo.UseDefault = true
		*bo.DefOpts = nats.GetDefaultOptions()
	} else {
		bo.UseDefault = false
		bo.CustomedOptFunc = opts
	}
}
