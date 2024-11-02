package main

import (
	"encoding/json"
	"github.com/EraldCaka/broker/pkg/config"
	kafkabroker "github.com/EraldCaka/broker/pkg/kafka"
	"github.com/EraldCaka/broker/services/users/types"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
)

func main() {

	config.InitConfig()
	subscriber, publisher := kafkabroker.NewClient()

	app := fiber.New()

	app.Get("/get", func(ctx *fiber.Ctx) error {
		log.Println("GET request triggered!")
		log.Println("IP:", ctx.IP())

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

	go func() {
		var targetStruct users.UserBody
		subscriber.ConsumeMessages(func(msg *message.Message) error {

			err := json.Unmarshal(msg.Payload, &targetStruct)
			if err != nil {
				log.Printf("Failed to unmarshal Kafka message: %v", err)
				return err
			}

			log.Printf("Kafka message received: %+v", targetStruct)

			return nil
		})
	}()

	log.Fatal(app.Listen(":5001"))
}
