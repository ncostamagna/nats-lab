package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
)

type Message struct {
	Content string `json:"content"`
}

func main() {
	// Connect to the NATS server in Minikube

	godotenv.Load()
	natsURL := os.Getenv("NATS_URL")
	credsPath := os.Getenv("NATS_CREDS_PATH")

	log.Println("NATS_URL: ", natsURL)
	log.Println("NATS_CREDS_PATH: ", credsPath)

	nc, err := nats.Connect(natsURL, nats.UserCredentials(credsPath))
	log.Println("NATS_URL: ", natsURL)

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
