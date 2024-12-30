# Broker

Broker will connect with configured services and optimally distribute messages
to from publishing services to receiving ones just like ```Nginx```. Using event-driven-architecture
and currently using round robin algorithm for the service discovery and selection.

## Methods for Service and Broker Interaction

### Initialization

```go
config.InitConfig()
consumer, publisher := kafkabroker.NewClient()
```

### Publish

Services publish messages to the Broker using the following method:
```go
func (p *Publisher) PublishMessage(topic string, msg []byte) error {
    return p.publisher.Publish(topic, message.NewMessage(watermill.NewUUID(), msg))
}
```

- Parameters:
  - ``topic``: The Kafka topic where the message should be published.
  - ``msg``: The payload to be sent.
- Usage: Services use the ``Publisher`` to publish messages to the corresponding topic. For example, in ``user-service``:

```go
err = publisher.PublishMessage(config.Config.Kafka.Services["user-service"].Topic, msg)
```

### Consume

Services consume messages from topics using the Subscriber:
```go
func (s *Subscriber) ConsumeMessages(routeFunc func(*message.Message) error) {
    for _, service := range config.Config.Kafka.Services {
        messages, err := s.subscriber.Subscribe(context.Background(), service.Topic)
        if err != nil {
            log.Fatalf("Failed to subscribe to topic: %v", err)
        }
        for msg := range messages {
            err = routeFunc(msg)
            if err != nil {
                log.Printf("Failed to route message: %v", err)
            }
            msg.Ack()
        }
    }
}
```

   - Parameters:
      - ``routeFunc``: A function to process each message. This is where custom logic is applied.
   - Usage: Services consume messages using this method and process them accordingly:

```go
func routeFunc(msg *message.Message) error {
	// error handling
	return nil
}

log.Printf("Received message with ID: %s", msg.UUID)

```
