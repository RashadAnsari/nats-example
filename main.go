package main

import (
	"log"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
)

type NatsConnection struct {
	Subject    string
	Queue      string
	Connection *nats.Conn
}

func (n *NatsConnection) Create() {
	var err error
	n.Connection, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("failed to create nates connection: %s", err.Error())
	}
}

func (n NatsConnection) Subscribe(name string) {
	_, err := n.Connection.QueueSubscribe(n.Subject, n.Queue, func(msg *nats.Msg) {
		log.Printf("receive message in subscriber %s: %s", name, string(msg.Data))
	})
	if err != nil {
		log.Fatalf("failed to create subsciber %s: %s", name, err.Error())
	}
}

func (n NatsConnection) Publish(value string) {
	if err := n.Connection.Publish(n.Subject, []byte(value)); err != nil {
		log.Fatalf("failed to publish: %s", err.Error())
	}
}

func main() {
	nc := NatsConnection{
		Subject: "my-subject",
		Queue:   "my-queue",
	}

	nc.Create()

	for i := 0; i < 2; i++ {
		nc.Subscribe(strconv.Itoa(i))
	}

	for i := 0; i < 1000; i++ {
		nc.Publish(strconv.Itoa(i))
	}

	time.Sleep(10 * time.Second)
}
