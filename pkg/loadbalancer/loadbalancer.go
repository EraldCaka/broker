package loadbalancer

import (
	"math/rand"
	"sync"
)

type LoadBalancer struct {
	services []string
	mu       sync.RWMutex
}

func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{}
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
	// in here I use RoundRobin logic
	return lb.services[rand.Intn(len(lb.services))]
}
