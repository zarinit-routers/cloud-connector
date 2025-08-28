package queue

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/charmbracelet/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zarinit-routers/cloud-connector/globals"
	"github.com/zarinit-routers/cloud-connector/models"
)

var (
	conn           *amqp.Connection
	requestsQueue  amqp.Queue
	responsesQueue amqp.Queue
	channel        *amqp.Channel
)

func Serve() {
	if connection, err := amqp.Dial("amqp://guest:guest@rabbit-mq:5672/"); err != nil {
		log.Fatal("Failed to connect to RabbitMQ", "error", err)
	} else {
		conn = connection
	}
	log.Info("Connected to RabbitMQ")
	defer conn.Close()

	if ch, err := conn.Channel(); err != nil {
		log.Fatal("Failed to open a channel", "error", err)
	} else {
		globals.SetChannel(ch)
	}

	messages, err := channel.Consume(
		requestsQueue.Name, // queue
		"",                 // consumer
		true,               // auto-ack
		false,              // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // args
	)
	if err != nil {
		log.Fatal("Failed to register a consumer", "error", err)
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for m := range messages {
			requestId := m.CorrelationId
			log.Info("Received a message", "message", string(m.Body), "requestId", requestId) // TODO: remove this log
			var cloudRequest models.FromCloudRequest
			if err := json.Unmarshal(m.Body, &cloudRequest); err != nil {
				log.Error("Failed to unmarshal message", "error", err)
				sendError(m.CorrelationId, BadRequestBodyErr(err))
				continue
			}
			if err := globals.SendRequest(cloudRequest.NodeID, cloudRequest.ToNode(requestId)); err != nil {
				log.Error("Failed to send request", "error", err)
				sendError(m.CorrelationId, err)
				continue
			}
		}
	}()

	wg.Wait()
}

func BadRequestBodyErr(err error) error {
	return fmt.Errorf("bad request body: %s", err)
}

func sendError(requestId string, err error) error {
	log.Error("Sending error response", "error", err, "requestId", requestId)
	response := &models.ToCloudResponse{
		RequestError: err.Error(),
	}

	return globals.SendResponse(requestId, response)
}
