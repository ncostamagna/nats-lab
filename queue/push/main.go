package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
)

func generateRandomNumber(max int64) int64 {
	n, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		log.Fatal(err)
	}
	return n.Int64() + 1 // Add 1 to make range 1 to max (inclusive)
}

func main() {

	godotenv.Load()
	natsURL := os.Getenv("NATS_HOST")
	seed := os.Getenv("NATS_SEED")
	kp, _ := nkeys.FromSeed([]byte(seed))
	pub, _ := kp.PublicKey()
	log.Println("NATS_HOST: ", natsURL)
	log.Println("NATS_SEED: ", seed)
	log.Println("kp: ", kp)

	nc, err := nats.Connect(natsURL, nats.Nkey(pub, func(data []byte) ([]byte, error) {
		return kp.Sign(data)
	}))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("NATS_HOST: ", natsURL)

	defer nc.Close()

	randomNum := generateRandomNumber(999999999)

	if err := nc.Publish("updates", []byte(fmt.Sprintf("%d", randomNum))); err != nil {
		log.Fatal(err)
	}

}
