package nats

import (
	"log"
	"testing"
)

func TestMsg(t *testing.T) {
	var msg *Msg
	msg.Subject = "demon"
	msg.Data = []byte("hello demon")

	res, err := msg.Encode()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(res))
}
