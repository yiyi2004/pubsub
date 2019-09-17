package nats

import (
	"os"
)

// NATS
var (
	NATSURL            = os.Getenv("NATS_SERVER_URL")
	DEFAULTMAXMESSAGES = os.Getenv("NATS_DEFAULT_MAX_MESSAGES")
)
