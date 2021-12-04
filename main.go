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
	MessageCount    = 10
	PublishDelay    = 100 * time.Millisecond
)

type NatsConnection struct {
	Subject string
	Queue   string
	Chan    chan string
	Conn    *nats.Conn
}

func (c NatsConnection) Subscribe(name string) {
	_, err := c.Conn.QueueSubscribe(c.Subject, c.Queue, func(msg *nats.Msg) {
		c.Chan <- fmt.Sprintf("receive message in nats subscriber %s: %s", name, string(msg.Data))
	})

	if err != nil {
		log.Fatalf("failed to create nats subsciber %s: %s", name, err.Error())
	}
}

func (c NatsConnection) Publish(value string) {
	if err := c.Conn.Publish(c.Subject, []byte(value)); err != nil {
		log.Fatalf("failed to publish to nats: %s", err.Error())
	}
}

type JetStreamConnection struct {
	Stream    string
	Subject   string
	Queue     string
	Chan      chan string
	Conn      *nats.Conn
	jsContext nats.JetStreamContext
}

func (c *JetStreamConnection) Init() {
	var err error

	c.jsContext, err = c.Conn.JetStream()
	if err != nil {
		log.Fatalf("failed to initialize jetstream connection: %s", err.Error())
	}

	_, err = c.jsContext.AddStream(&nats.StreamConfig{
		Name:     c.Stream,
		Subjects: []string{c.Subject},
	})

	if err != nil {
		log.Fatalf("failed to add stream: %s", err.Error())
	}
}

func (c JetStreamConnection) Subscribe(name string) {
	opts := []nats.SubOpt{
		nats.Durable(c.Queue),
		nats.AckAll(),
		nats.ManualAck(),
		nats.DeliverLast(),
	}

	_, err := c.jsContext.QueueSubscribe(c.Subject, c.Queue, func(msg *nats.Msg) {
		c.Chan <- fmt.Sprintf("receive message in jetstream subscriber %s: %s", name, string(msg.Data))

		if err := msg.Ack(); err != nil {
			log.Fatalf("failed to send ack for jetstream message %s: %s", name, string(msg.Data))
		}
	}, opts...)

	if err != nil {
		log.Fatalf("failed to create jestream subsciber %s: %s", name, err.Error())
	}
}

func (c JetStreamConnection) Publish(value string) {
	if _, err := c.jsContext.Publish(c.Subject, []byte(value)); err != nil {
		log.Fatalf("failed to publish to jetstream: %s", err.Error())
	}
}

func main() {
	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("failed to create nates connection: %s", err.Error())
	}

	nc := NatsConnection{
		Subject: "my-subject-nats",
		Queue:   "my-queue-nats",
		Chan:    make(chan string),
		Conn:    conn,
	}

	for i := 1; i <= SubscriberCount; i++ {
		nc.Subscribe(strconv.Itoa(i))
	}

	go func() {
		for i := 1; i <= MessageCount; i++ {
			nc.Publish(strconv.Itoa(i))
			time.Sleep(PublishDelay)
		}
	}()

	for i := 1; i <= MessageCount; i++ {
		msg := <-nc.Chan
		log.Println(msg)
	}

	fmt.Println("==========")

	js := JetStreamConnection{
		Stream:  "my-stream-jetstream",
		Subject: "my-subject-jetstream",
		Queue:   "my-queue-jetstream",
		Chan:    make(chan string),
		Conn:    conn,
	}

	js.Init()

	for i := 1; i <= SubscriberCount; i++ {
		js.Subscribe(strconv.Itoa(i))
	}

	go func() {
		for i := 1; i <= MessageCount; i++ {
			js.Publish(strconv.Itoa(i))
			time.Sleep(PublishDelay)
		}
	}()

	for i := 1; i <= MessageCount; i++ {
		msg := <-js.Chan
		log.Println(msg)
	}
}
