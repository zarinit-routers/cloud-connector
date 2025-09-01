package main

import (
	"flag"
	"fmt"

	"github.com/zarinit-routers/cloud-connector/connections"
)

var (
	printCount = flag.Uint("count", 1, "Number of nodes to generate")
)

func init() {
	flag.Parse()
}

func main() {

	for i := 0; i < int(*printCount); i++ {
		printNode()
	}
}

func printNode() {
	fmt.Println(connections.GenNodeName())
}
