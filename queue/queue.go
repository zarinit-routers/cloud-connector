package queue

import (
	"github.com/charmbracelet/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Serve() {
	conn, err := amqp.Dial("amqp://guest:guest@rabbit-mq:5672/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ", "error", err)
	}
	log.Info("Connected to RabbitMQ")
	defer conn.Close()
}
