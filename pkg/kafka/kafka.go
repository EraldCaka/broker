package kafka

import (
	"context"
	"github.com/EraldCaka/broker/pkg/config"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"log"
)

type Subscriber struct {
	subscriber *kafka.Subscriber
}

type Publisher struct {
	publisher *kafka.Publisher
}

func NewClient() (*Subscriber, *Publisher) {
	kafkaConfig := kafka.SubscriberConfig{
		Brokers:       []string{config.Get("kafka.bootstrap_servers")},
		Unmarshaler:   kafka.DefaultMarshaler{},
		ConsumerGroup: "example-group",
	}

	kafkaSub, err := kafka.NewSubscriber(kafkaConfig, watermill.NewStdLogger(false, false))
	if err != nil {
		log.Fatalf("Failed to create Kafka subscriber: %v", err)
	}

	kafkaPubConfig := kafka.PublisherConfig{
		Brokers:   []string{config.Get("kafka.bootstrap_servers")},
		Marshaler: kafka.DefaultMarshaler{},
	}

	kafkaPub, err := kafka.NewPublisher(kafkaPubConfig, watermill.NewStdLogger(false, false))
	if err != nil {
		log.Fatalf("Failed to create Kafka publisher: %v", err)
	}

	return &Subscriber{subscriber: kafkaSub}, &Publisher{publisher: kafkaPub}
}

func (s *Subscriber) ConsumeMessages(routeFunc func(*message.Message) error) {
	messages, err := s.subscriber.Subscribe(context.Background(), config.Get("kafka.topic"))
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %v", err)
	}

	for msg := range messages {
		err := routeFunc(msg)
		if err != nil {
			log.Printf("Failed to route message: %v", err)
		}
		msg.Ack()
	}
}

func (p *Publisher) PublishMessage(topic string, msg []byte) error {
	return p.publisher.Publish(topic, message.NewMessage(watermill.NewUUID(), msg))
}
