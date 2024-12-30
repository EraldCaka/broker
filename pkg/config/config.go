package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type ServiceConfig struct {
	Topic        string   `yaml:"topic"`
	Role         string   `yaml:"role"`
	Url          string   `yaml:"url"`
	Name         string   `yaml:"name"`
	SubscribesTo []string `yaml:"subscribes_to"`
}

type KafkaConfig struct {
	BootstrapServers string                   `yaml:"bootstrap_servers"`
	Services         map[string]ServiceConfig `yaml:"services"`
}

type AppConfig struct {
	Kafka KafkaConfig `yaml:"kafka"`
}

var Config AppConfig

func InitConfig() {
	data, err := os.ReadFile("/Users/eraldcaka/Documents/github/broker/config.yaml")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	err = yaml.Unmarshal(data, &Config)
	if err != nil {
		log.Fatalf("Failed to unmarshal config data: %v", err)
	}
}
