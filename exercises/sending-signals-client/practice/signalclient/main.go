package main

import (
	"context"
	"log"
	"signals"

	// TODO Part D: Add "context" to your module imports.

	"go.temporal.io/sdk/client"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	signal := signals.FulfillOrderSignal{Fulfilled: true}
	err = c.SignalWorkflow(
		context.Background(),
		"signals",
		"",
		"fulfill-order-signal",
		signal,
	)
	if err != nil {
		log.Fatalln("Error sending the Signal", err)
		return
	}
}
