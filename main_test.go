package main

import (
	"testing"

	"github.com/nats-io/nats.go"
	"time"
)

func TestPublishAndSubscribe(t *testing.T) {
	nc, ns, err := RunEmbeddedServer(true, false)
	if err != nil {
		t.Fatal(err)
	}
	defer ns.Shutdown()

	// Create a channel to receive the message
	msgChan := make(chan *nats.Msg, 1)

	// Subscribe to messages
	sub, err := nc.Subscribe("hello.world", func(msg *nats.Msg) {
		msgChan <- msg
		msg.Respond([]byte("Ahoy there!"))
	})
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Unsubscribe()

	// Publish a message
	err = PublishMessage(nc, "hello.world", "Test message")
	if err != nil {
		t.Fatal(err)
	}

	// Wait for the message with a timeout
	select {
	case msg := <-msgChan:
		if string(msg.Data) != "Test message" {
			t.Errorf("Expected message 'Test message', got '%s'", string(msg.Data))
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for message")
	}
}

func BenchmarkRequestReplyLoopback(b *testing.B) {
	nc, ns, err := RunEmbeddedServer(false, false)
	if err != nil {
		b.Fatal(err)
	}

	nc.Subscribe("hello.world", func(msg *nats.Msg) {
		msg.Respond([]byte("Hi there"))
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := nc.Request("hello.world", []byte("hihi"), 10*time.Second)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.StopTimer()
	ns.Shutdown()
	ns.WaitForShutdown()
}

func BenchmarkRequestReplyInProcess(b *testing.B) {
	nc, ns, err := RunEmbeddedServer(true, false)
	if err != nil {
		b.Fatal(err)
	}

	nc.Subscribe("hello.world", func(msg *nats.Msg) {
		msg.Respond([]byte("Hi there"))
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := nc.Request("hello.world", []byte("hihi"), 10*time.Second)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.StopTimer()
	ns.Shutdown()
	ns.WaitForShutdown()
}