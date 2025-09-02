package main

import (
	"fmt"
	"go-microservices/message-service/rabbitmq"
	initializers "go-microservices/pkg/initializer"
	"log"
)

func init() {
	initializers.LoadEnvvariables()
}
func main() {
	conn, ch, msgs := rabbitmq.InitRabbitMQ()
	defer conn.Close()
	defer ch.Close()

	// Start consumer
	go func() {
		for d := range msgs {
			// Pretend this is sending SMS
			fmt.Printf("📩 [SMS Service] Sending SMS for new user: %s\n", d.Body)
		}
	}()

	log.Println("🚀 SMS service is running. Waiting for messages...")
	select {}
}