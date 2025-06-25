package main

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
)

func main() {

	godotenv.Load()
	natsURL := os.Getenv("NATS_HOST")
	seed := os.Getenv("NATS_SEED")
	kp, _ := nkeys.FromSeed([]byte(seed))
	pub, _ := kp.PublicKey()

	nc, err := nats.Connect(natsURL, nats.Nkey(pub, func(data []byte) ([]byte, error) {
		return kp.Sign(data)
	}))
	if err != nil {
		log.Fatal(err)
	}

	defer nc.Close()

	// Use a WaitGroup to wait for 10 messages to arrive
	wg := sync.WaitGroup{}
	wg.Add(10)

	// Create a queue subscription on "updates" with queue name "workers"
	if _, err := nc.QueueSubscribe("updates", "workers", func(m *nats.Msg) {
		log.Println("Queue sub #1: ", string(m.Data))
		wg.Done()
	}); err != nil {
		log.Fatal(err)
	}

	if _, err := nc.QueueSubscribe("updates", "workers", func(m *nats.Msg) {
		log.Println("Queue sub #2: ", string(m.Data))
		wg.Done()
	}); err != nil {
		log.Fatal(err)
	}

	if _, err := nc.Subscribe("updates", func(m *nats.Msg) {
		log.Println("Normal sub #1: ", string(m.Data))
		wg.Done()
	}); err != nil {
		log.Fatal(err)
	}

	if _, err := nc.Subscribe("updates", func(m *nats.Msg) {
		log.Println("Normal sub #2: ", string(m.Data))
		wg.Done()
	}); err != nil {
		log.Fatal(err)
	}

	// Wait for messages to come in
	wg.Wait()

}
