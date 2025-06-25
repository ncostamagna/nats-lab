package main

import (
	"fmt"
	"log"
	"os"

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

	js, _ := nc.JetStream()

	

	kv, err := js.KeyValue("TEST")
	if err != nil {
		fmt.Println("Create bucket ")
		kv, err = js.CreateKeyValue(&nats.KeyValueConfig{
			Bucket: "TEST",
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	kw, err := kv.Watch("key")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for e := range kw.Updates() {
			if e != nil {
				fmt.Println(e)
				fmt.Println(e.Operation())
				fmt.Println(e.Revision())
				fmt.Println("kv: ", e.Key, string(e.Value()))
			}
		}
	}()

	v, _ := kv.Get("key")
	fmt.Println("kv: ", v)
	kv.Put("key", []byte("value"))

	v, _ = kv.Get("key")
	fmt.Println("kv: ", string(v.Value()))


	kv.Update("key", []byte("value2"), v.Revision())

	v, _ = kv.Get("key")
	fmt.Println("kv: ", string(v.Value()))


	kv.Delete("key")

	v, _ = kv.Get("key")
	fmt.Println("kv: ", v)


	
	
}