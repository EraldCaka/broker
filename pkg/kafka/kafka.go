package kafka

import (
	"context"
	"log"

	"github.com/EraldCaka/broker/pkg/config"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
)

type Subscriber struct {
	subscriber *kafka.Subscriber
}

type Publisher struct {
	publisher *kafka.Publisher
}

func NewClient() (*Subscriber, *Publisher) {

	kafkaConfig := kafka.SubscriberConfig{
		Brokers:       []string{config.Config.Kafka.BootstrapServers},
		Unmarshaler:   kafka.DefaultMarshaler{},
		ConsumerGroup: "example-group",
	}

	kafkaSub, err := kafka.NewSubscriber(kafkaConfig, watermill.NewStdLogger(false, false))
	if err != nil {
		log.Fatalf("Failed to create Kafka subscriber: %v", err)
	}
	log.Println("Kafka subscriber created successfully")
	kafkaPubConfig := kafka.PublisherConfig{
		Brokers:   []string{config.Config.Kafka.BootstrapServers},
		Marshaler: kafka.DefaultMarshaler{},
	}

	kafkaPub, err := kafka.NewPublisher(kafkaPubConfig, watermill.NewStdLogger(false, false))
	if err != nil {
		log.Fatalf("Failed to create Kafka publisher: %v", err)
	}
	log.Println("Kafka publisher created successfully")
	return &Subscriber{subscriber: kafkaSub}, &Publisher{publisher: kafkaPub}
}

func (s *Subscriber) ConsumeMessages(routeFunc func(*message.Message) error) {
	for _, service := range config.Config.Kafka.Services {
		subscriptionTopics := service.SubscribesTo
		for _, topic := range subscriptionTopics {
			log.Printf("Starting subscription for service %s on topic %s", service.Name, topic)
			messages, err := s.subscriber.Subscribe(context.Background(), topic)
			if err != nil {
				log.Fatalf("Failed to subscribe to topic %s: %v", topic, err)
			}
			go func(service config.ServiceConfig, topic string) {
				for msg := range messages {
					log.Printf("Received message on topic %s for service %s: %s", topic, service.Name, string(msg.Payload))
					err := routeFunc(msg)
					if err != nil {
						log.Printf("Failed to route message: %v", err)
					}
					msg.Ack()
				}
			}(service, topic)
		}
	}
}

func (p *Publisher) PublishMessage(topic string, msg []byte) error {
	return p.publisher.Publish(topic, message.NewMessage(watermill.NewUUID(), msg))
}
