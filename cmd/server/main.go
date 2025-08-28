package main

import (
	"sync"

	"github.com/zarinit-routers/cloud-connector/connections"
	"github.com/zarinit-routers/cloud-connector/queue"
)

func main() {
	wg := sync.WaitGroup{}

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

	wg.Wait()
}
