package main

import (
	"github.com/EraldCaka/broker/pkg/broker"
	"github.com/EraldCaka/broker/pkg/config"
	"log"
)

func main() {

	config.InitConfig()

	b := broker.NewBroker()
	if err := b.Start(); err != nil {
		log.Fatalf("Failed to start broker: %v", err)
	}
}
