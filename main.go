package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	SubscriberCount = 2
	MessageCount    = 100
	PublishDelay    = 100 * time.Millisecond
)

type NatsConnection struct {
	Subject    string
	Queue      string
	Chan       chan string
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
		n.Chan <- fmt.Sprintf("receive message in subscriber %s: %s", name, string(msg.Data))
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
		Chan:    make(chan string),
	}

	nc.Create()

	for i := 0; i < SubscriberCount; i++ {
		nc.Subscribe(strconv.Itoa(i))
	}

	go func() {
		for i := 0; i < MessageCount; i++ {
			nc.Publish(strconv.Itoa(i))
			time.Sleep(PublishDelay)
		}
	}()

	for i := 0; i < MessageCount; i++ {
		msg := <-nc.Chan
		log.Println(msg)
	}
}
