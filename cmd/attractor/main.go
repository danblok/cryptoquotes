package main

import (
	"log"
	"os"

	"github.com/danblok/cryptoquotes/internal/messages"
	"github.com/danblok/cryptoquotes/pkg/config"
)

func main() {
	m := messages.New(config.NewAMQPConfig(os.Getenv("MQ_LOGIN"), os.Getenv("MQ_PASSWORD"), "mq", 5672, nil))
	_, _, err := m.ReceivivgChannelAndQueue()
	if err != nil {
		log.Fatal(err)
	}
	defer m.Close()

	if err := m.LogMessages(); err != nil {
		log.Fatal(err)
	}
}
