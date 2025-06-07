package main

import (
	"context"
	"encoding/base64"
	"flag"
	"log"

	"go.temporal.io/sdk/client"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	var taskToken string
	flag.StringVar(&taskToken, "tasktoken", "", "Task Token of Async. Activity to Complete")
	flag.Parse()
	// Decode from hex…
	//decoded, err := hex.DecodeString(taskToken)
	// Or Base64
	decoded, err := base64.StdEncoding.DecodeString(taskToken)
	if err != nil {
		log.Fatalln("Unable to decode token", err)
	}

	var result string
	err = c.CompleteActivity(context.Background(), decoded, result, err)
	if err != nil {
		log.Fatalln("Unable to complete Async. Activity")
	}
}
