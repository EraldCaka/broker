package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type ServiceConfig struct {
	Topic        string   `yaml:"topic"`
	Role         string   `yaml:"role"`
	SubscribesTo []string `yaml:"subscribes_to"`
}

type KafkaConfig struct {
	BootstrapServers string                   `yaml:"bootstrap_servers"`
	Services         map[string]ServiceConfig `yaml:"services"`
}

type AppConfig struct {
	Kafka KafkaConfig `yaml:"kafka"`
}

// Config Global variable to hold the configuration
var Config AppConfig

func InitConfig() {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	err = yaml.Unmarshal(data, &Config)
	if err != nil {
		log.Fatalf("Failed to unmarshal config data: %v", err)
	}
}
