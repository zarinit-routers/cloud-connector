package main

import (
	"encoding/json"
	"sync"

	"github.com/charmbracelet/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zarinit-routers/cloud-connector/connections"
	"github.com/zarinit-routers/cloud-connector/models"
	"github.com/zarinit-routers/cloud-connector/queue"
	"github.com/zarinit-routers/cloud-connector/server"
	"github.com/zarinit-routers/cloud-connector/storage/database"
)

func main() {
	wg := sync.WaitGroup{}

	if err := database.Setup(); err != nil {
		log.Fatal("Failed to setup database", "error", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		connections.Serve()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		queue.Serve()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.Serve(); err != nil {
			log.Fatal("Failed serve HTTP server", "error", err)
		}
	}()

	queue.AddHandler(queueHandler)

	connections.AddHandler(websocketHandler)

	wg.Wait()
}

func queueHandler(m *amqp.Delivery) error {
	requestId := m.CorrelationId
	var cloudRequest models.FromCloudRequest
	if err := json.Unmarshal(m.Body, &cloudRequest); err != nil {
		log.Error("Failed to unmarshal message", "error", err)
		return err
	}
	if err := connections.SendRequest(cloudRequest.NodeID, cloudRequest.ToNode(requestId)); err != nil {
		log.Error("Failed to send request", "error", err)
		return err
	}
	return nil
}

func websocketHandler(body []byte) error {
	var request models.FromNodeResponse
	if err := json.Unmarshal(body, &request); err != nil {
		log.Error("Failed to unmarshal message", "error", err)
		return err
	}
	if err := queue.SendResponse(request.RequestID, request.ToCloud()); err != nil {
		log.Error("Failed to send response", "error", err)
		return err
	}
	return nil
}
