package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

const (
	SubscriberCount = 2
	MessageCount    = 100
	PublishDelay    = 100 * time.Millisecond

	ClusterName = "test-cluster"
	ClientID    = "test-123"
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
		n.Chan <- fmt.Sprintf("receive message in nats subscriber %s: %s", name, string(msg.Data))
	})
	if err != nil {
		log.Fatalf("failed to create nats subsciber %s: %s", name, err.Error())
	}
}

func (n NatsConnection) Publish(value string) {
	if err := n.Connection.Publish(n.Subject, []byte(value)); err != nil {
		log.Fatalf("failed to publish to nats: %s", err.Error())
	}
}

type StanConnection struct {
	Subject    string
	Queue      string
	Chan       chan string
	Connection stan.Conn
}

func (s *StanConnection) Create() {
	var err error
	s.Connection, err = stan.Connect(ClusterName, ClientID, stan.NatsURL(nats.DefaultURL))

	if err != nil {
		log.Fatalf("failed to create nates connection: %s", err.Error())
	}
}

func (s StanConnection) Subscribe(name string) {
	_, err := s.Connection.QueueSubscribe(s.Subject, s.Queue, func(msg *stan.Msg) {
		s.Chan <- fmt.Sprintf("receive message in stan subscriber %s: %s", name, string(msg.Data))

		if err := msg.Ack(); err != nil {
			log.Fatalf("failed to send ack to stan in subscriber %s: %s", name, err.Error())
		}
	}, stan.SetManualAckMode())
	if err != nil {
		log.Fatalf("failed to create stan subsciber %s: %s", name, err.Error())
	}
}

func (s StanConnection) Publish(value string) {
	if err := s.Connection.Publish(s.Subject, []byte(value)); err != nil {
		log.Fatalf("failed to publish to stan: %s", err.Error())
	}
}

func main() {
	nc := NatsConnection{
		Subject: "my-subject-nats",
		Queue:   "my-queue-nats",
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

	log.Println("==================================================")

	sc := StanConnection{
		Subject: "my-subject-stan",
		Queue:   "my-queue-stan",
		Chan:    make(chan string),
	}

	sc.Create()

	for i := 0; i < SubscriberCount; i++ {
		sc.Subscribe(strconv.Itoa(i))
	}

	go func() {
		for i := 0; i < MessageCount; i++ {
			sc.Publish(strconv.Itoa(i))
			time.Sleep(PublishDelay)
		}
	}()

	for i := 0; i < MessageCount; i++ {
		msg := <-sc.Chan
		log.Println(msg)
	}
}
