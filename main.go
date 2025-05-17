package main

import (
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	// Connect to the NATS server in Minikube
	nc, err := nats.Connect("nats://127.0.0.1:4222")  // Default Minikube IP
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	sub, err := nc.Subscribe("hello.world", func(msg *nats.Msg) {
		log.Printf("Received message: %s", string(msg.Data))
		msg.Respond([]byte("Ahoy there!"))
	})
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()

	err = PublishMessage(nc, "hello.world", "Hello from Go!")
	if err != nil {
		log.Printf("Error publishing message: %v", err)
	}

	log.Println("Connected to NATS server. Press Ctrl+C to exit.")
	
	// Keep the program running
	select {}
}

// PublishMessage publishes a message to a NATS subject
func PublishMessage(nc *nats.Conn, subject string, message string) error {
	return nc.Publish(subject, []byte(message))
}
