package rabbitmq

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

var RabbitConn *amqp.Connection
var RabbitCh *amqp.Channel

func SetupRabbitMQ() {
    var err error
	url := os.Getenv("RABBITMQ_URL")
    RabbitConn, err = amqp.Dial(url)
    if err != nil {
        log.Fatalf("Failed to connect to RabbitMQ: %v", err)
    }

    RabbitCh, err = RabbitConn.Channel()
    if err != nil {
        log.Fatalf("Failed to open RabbitMQ channel: %v", err)
    }

    // Declare a queue (idempotent: safe to call multiple times)
    _, err = RabbitCh.QueueDeclare(
        "user.created", // queue name
        true,           // durable
        false,          // autoDelete
        false,          // exclusive
        false,          // noWait
        nil,            // args
    )
    if err != nil {
        log.Fatalf("Failed to declare queue: %v", err)
    }
}