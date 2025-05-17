package main

import (
	"log"
	"net/http"
	"encoding/json"

	"github.com/nats-io/nats.go"
)

type Message struct {
	Content string `json:"content"`
}

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
	
	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var message Message
		if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Publish message to NATS
		if err := PublishMessage(nc, "hello.world", message.Content); err != nil {
			http.Error(w, "Failed to publish message", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Message sent successfully"))
	})

	// Start HTTP server
	log.Println("Starting HTTP server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// PublishMessage publishes a message to a NATS subject
func PublishMessage(nc *nats.Conn, subject string, message string) error {
	return nc.Publish(subject, []byte(message))
}
