# NATS Example

An example in Go of two ways to share work between subscribers with [NATS](https://nats.io/): a core NATS queue group, and a durable JetStream queue consumer.

## What it does

[main.go](main.go) runs two rounds, each with 2 subscribers in the same queue group and 10 messages published 100 ms apart:

1. **Core NATS**: `QueueSubscribe` on `my-subject-nats`. Each message goes to exactly one subscriber in the group. Nothing is stored, so a message published while no one is listening is lost.
2. **JetStream**: creates the stream `my-stream-jetstream` on `my-subject-jetstream` and subscribes with a durable queue consumer (`Durable`, `ManualAck`, `AckAll`, `DeliverLast`). Messages are persisted by the server and each one is acknowledged after it is handled.

Each received message is logged with the name of the subscriber that got it, so you can see the load being split.

## Run

Requires Go and Docker.

```bash
make up      # start nats-server with JetStream (-js) on localhost:4222
go run .
make down
```

Example output:

```
receive message in nats subscriber 1: 1
receive message in nats subscriber 2: 2
...
==========
receive message in jetstream subscriber 2: 1
...
```

## License

[MIT](LICENSE)
