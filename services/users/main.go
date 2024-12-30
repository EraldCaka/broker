package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/EraldCaka/broker/pkg/config"
	kafkabroker "github.com/EraldCaka/broker/pkg/kafka"
	users "github.com/EraldCaka/broker/services/users/types"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gofiber/fiber/v2"
)

func routeFunc(msg *message.Message) error {
	// error handling
	return nil
}

/*
dummy service for testing purposes only
*/
func main() {
	config.InitConfig()
	consumer, publisher := kafkabroker.NewClient()

	app := fiber.New()

	app.Get("/get", func(ctx *fiber.Ctx) error {
		log.Println("GET request triggered!")

		user := &users.UserBody{
			Name:     "user1",
			Lastname: "temporary",
			Age:      35,
		}

		msg, err := json.Marshal(user)
		if err != nil {
			log.Printf("Failed to serialize user data: %v", err)
			return ctx.Status(http.StatusInternalServerError).SendString("Failed to process request")
		}

		err = publisher.PublishMessage(config.Config.Kafka.Services["user-service"].Topic, msg)
		if err != nil {
			log.Printf("Failed to publish message to Kafka: %v", err)
			return ctx.Status(http.StatusInternalServerError).SendString("Failed to process request")
		}

		log.Println("User data sent to Kafka:", user)
		return ctx.JSON(user)
	})

	go consumer.ConsumeMessages(routeFunc)

	log.Fatal(app.Listen(":5001"))
}
