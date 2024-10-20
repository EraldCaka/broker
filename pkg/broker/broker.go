package broker

import (
	"github.com/EraldCaka/broker/pkg/discovery"
	"github.com/EraldCaka/broker/pkg/kafka"
	"github.com/EraldCaka/broker/pkg/loadbalancer"
	"github.com/ThreeDotsLabs/watermill/message"
	"log"
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

	go b.kafkaSubscriber.ConsumeMessages(b.routeMessage)
	go b.serviceRegistry.RegisterService("auth-service", "localhost:5000")
	go b.serviceRegistry.RegisterService("user-service", "localhost:5001")
	go b.serviceRegistry.RegisterService("product-service", "localhost:5002")

	select {}
}

func (b *Broker) routeMessage(msg *message.Message) error {
	service := b.loadBalancer.SelectService()
	log.Printf("Routing message %s to service %s", string(msg.Payload), service)

	return nil
}
