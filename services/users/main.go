// user service example
package main

import (
	"github.com/EraldCaka/broker/pkg/config"
	kafkabroker "github.com/EraldCaka/broker/pkg/kafka"
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
		log.Println("get request triggered!")
		log.Println("ip:", ctx.IP())

		msg := []byte("Request received on /get endpoint")
		err := publisher.PublishMessage(config.Config.Kafka.Services["user-service"].Topic, msg)
		if err != nil {
			log.Printf("Failed to publish message to Kafka: %v", err)
			return ctx.Status(http.StatusInternalServerError).SendString("Failed to process request")
		}

		return ctx.Status(http.StatusOK).JSON("user1")
	})

	go func() {
		subscriber.ConsumeMessages(func(msg *message.Message) error {
			log.Printf("Kafka message received: %s", string(msg.Payload))
			return nil
		})
	}()

	app.Listen(":5001")
}
