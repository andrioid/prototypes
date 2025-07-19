package main

import (
	"log"
	"nats-w-auth/pkg/nats_scs"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func setupSession() {
	nc, err := nats.Connect(natsd.ClientURL())
	if err != nil {
		log.Fatal("Failed to create nats client", err)
	}
	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal("Failed to create jetstream", err)
	}
	// Initialize a new session manager and configure the session lifetime.
	natsStore := nats_scs.New(js)
	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Store = natsStore
}
