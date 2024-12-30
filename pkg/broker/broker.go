package broker

import (
	"fmt"
	"log"

	"github.com/EraldCaka/broker/pkg/config"
	"github.com/EraldCaka/broker/pkg/discovery"
	"github.com/EraldCaka/broker/pkg/kafka"
	"github.com/EraldCaka/broker/pkg/loadbalancer"
	"github.com/ThreeDotsLabs/watermill/message"
)

type Broker struct {
	kafkaSubscriber *kafka.Subscriber
	kafkaPublisher  *kafka.Publisher
	serviceRegistry *discovery.ServiceRegistry
	loadBalancer    *loadbalancer.LoadBalancer
}

func NewBroker() *Broker {
	kafkaSub, kafkaPub := kafka.NewClient()
	serviceRegistry := discovery.NewServiceRegistry()
	loadBalancer := loadbalancer.NewLoadBalancer()

	return &Broker{
		kafkaSubscriber: kafkaSub,
		kafkaPublisher:  kafkaPub,
		serviceRegistry: serviceRegistry,
		loadBalancer:    loadBalancer,
	}
}

func (b *Broker) Start() error {
	log.Println("Starting broker...")
	b.setServices()

	select {}
}

func (b *Broker) setServices() {

	go b.kafkaSubscriber.ConsumeMessages(b.routeMessage)

	for _, service := range config.Config.Kafka.Services {
		go b.serviceRegistry.RegisterService(service.Name, service.Url)
	}
}

func (b *Broker) routeMessage(msg *message.Message) error {
	targetService := b.loadBalancer.SelectService()

	if targetService == "" {
		return fmt.Errorf("no available services to route message")
	}

	var serviceConfig config.ServiceConfig
	for _, service := range config.Config.Kafka.Services {
		if service.Url == targetService {
			serviceConfig = service
			break
		}
	}

	log.Printf("Routing message with id %s to service %s (role: %s)", msg.UUID, serviceConfig.Name, serviceConfig.Role)

	switch serviceConfig.Role {
	case "publisher":
		return b.kafkaPublisher.PublishMessage(serviceConfig.Topic, msg.Payload)
	case "subscriber":
		log.Printf("Service %s is a subscriber; forwarding message to subscribed topics", serviceConfig.Name)
		for _, subscribedTopic := range serviceConfig.SubscribesTo {
			err := b.kafkaPublisher.PublishMessage(subscribedTopic, msg.Payload)
			if err != nil {
				log.Printf("Failed to forward message to subscribed topic %s: %v", subscribedTopic, err)
			}
		}
		return nil
	case "both":
		err := b.kafkaPublisher.PublishMessage(serviceConfig.Topic, msg.Payload)
		if err != nil {
			log.Printf("Failed to publish message to topic %s: %v", serviceConfig.Topic, err)
		}
		for _, subscribedTopic := range serviceConfig.SubscribesTo {
			err := b.kafkaPublisher.PublishMessage(subscribedTopic, msg.Payload)
			if err != nil {
				log.Printf("Failed to forward message to subscribed topic %s: %v", subscribedTopic, err)
			}
		}
		return nil
	default:
		return fmt.Errorf("unknown role for service %s", serviceConfig.Name)
	}
}
