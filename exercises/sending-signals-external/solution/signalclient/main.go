package main

import (
	"context"
	"log"

	"go.temporal.io/sdk/client"
)

type FulfillOrderSignal struct {
	Fulfilled bool
}

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	signal := FulfillOrderSignal{
		Fulfilled: true,
	}

	err = c.SignalWorkflow(context.Background(), "signals", "", "fulfill-order-signal", signal)
	if err != nil {
		log.Fatalln("Error sending the Signal", err)
		return
	}
}
