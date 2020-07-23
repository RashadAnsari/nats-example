package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"log"
	"strconv"
	"time"
)

const (
	ClusterName = "test-cluster"
	ClientID    = "test-123"
)

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
		s.Chan <- fmt.Sprintf("receive message in subscriber %s: %s", name, string(msg.Data))
	})
	if err != nil {
		log.Fatalf("failed to create subsciber %s: %s", name, err.Error())
	}
}

func (s StanConnection) Publish(value string) {
	if err := s.Connection.Publish(s.Subject, []byte(value)); err != nil {
		log.Fatalf("failed to publish: %s", err.Error())
	}
}

func main() {
	sc := StanConnection{
		Subject: "my-subject",
		Queue:   "my-queue",
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
