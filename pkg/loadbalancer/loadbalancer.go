package loadbalancer

import (
	"math/rand"
	"sync"

	"github.com/EraldCaka/broker/pkg/config"
)

type LoadBalancer struct {
	services []string
	mu       sync.RWMutex
}

func NewLoadBalancer() *LoadBalancer {
	var services []string
	for _, service := range config.Config.Kafka.Services {
		services = append(services, service.Url)
	}
	return &LoadBalancer{services: services}
}

func (lb *LoadBalancer) RegisterService(service string) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	lb.services = append(lb.services, service)
}

func (lb *LoadBalancer) DeregisterService(service string) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	for i, s := range lb.services {
		if s == service {
			lb.services = append(lb.services[:i], lb.services[i+1:]...)
			return
		}
	}
}

func (lb *LoadBalancer) SelectService() string {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	if len(lb.services) == 0 {
		return ""
	}
	return lb.services[rand.Intn(len(lb.services))]
}
