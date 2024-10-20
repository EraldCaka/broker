package discovery

import (
	"log"
	"sync"
)

type ServiceRegistry struct {
	services map[string]string
	mu       sync.RWMutex
}

func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		services: make(map[string]string),
	}
}

func (sr *ServiceRegistry) RegisterService(name, address string) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	sr.services[name] = address
	log.Printf("Service %s registered at %s", name, address)
}

func (sr *ServiceRegistry) DeregisterService(name string) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	delete(sr.services, name)
	log.Printf("Service %s deregistered", name)
}

func (sr *ServiceRegistry) GetService(name string) (string, bool) {
	sr.mu.RLock()
	defer sr.mu.RUnlock()
	address, ok := sr.services[name]
	return address, ok
}
