package globals

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zarinit-routers/cloud-connector/models"
)

var (
	connections = map[string]*websocket.Conn{}

	requestsQueue  amqp.Queue
	responsesQueue amqp.Queue
	channel        *amqp.Channel
)

func SetChannel(ch *amqp.Channel) {
	channel = ch

	setupQueues()
}

func GetRequestsFromQueue() (<-chan amqp.Delivery, error) {
	return channel.Consume(
		requestsQueue.Name, // queue
		"",                 // consumer
		true,               // auto-ack
		false,              // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // args
	)
}

func setupQueues() {

	q, err := channel.QueueDeclare(
		"requests", // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		log.Fatal("Failed to declare a queue", "error", err)
	} else {
		requestsQueue = q
	}

	q2, err := channel.QueueDeclare(
		"responses", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		log.Fatal("Failed to declare a queue", "error", err)
	} else {
		responsesQueue = q2
	}
}

func AppendConnection(nodeId string, conn *websocket.Conn) {
	connections[nodeId] = conn
}

func SendResponse(requestId string, response *models.ToCloudResponse) error {
	body, err := json.Marshal(response)
	if err != nil {
		return err
	}
	return channel.Publish(
		"",                  // exchange
		responsesQueue.Name, // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: requestId,
			Body:          body,
		},
	)
}

func SendRequest(nodeId string, r *models.ToNodeRequest) error {
	conn, ok := connections[nodeId]
	if !ok {
		return fmt.Errorf("node with id %q not connected", nodeId)
	}

	message, err := json.Marshal(r)
	if err != nil {
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, message)
}
