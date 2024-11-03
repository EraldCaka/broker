package broker

import (
	"github.com/EraldCaka/broker/pkg/config"
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
	service := b.loadBalancer.SelectService()
	log.Printf("Routing message with id %s to service %s", msg.UUID, service)

	// TODO: implement the service proxy
	return nil
}
