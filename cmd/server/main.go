package main

import (
	"sync"

	"github.com/zarinit-routers/cloud-connector/connections"
)

func main() {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		connections.Serve()
	}()

	wg.Wait()
}
