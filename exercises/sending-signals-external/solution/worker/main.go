package main

import (
	"log"

	interacting "edu-interacting-with-workflows-go-code/exercises/sending-signals-external/solution"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "signals", worker.Options{})

	w.RegisterWorkflow(interacting.Workflow)
	w.RegisterActivity(interacting.Activity)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
