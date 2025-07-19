package main

import (
	"fmt"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

func setupNats() {
	var err error
	opts := server.Options{
		JetStream: true,
		StoreDir:  "./data",
	}
	natsd, err = server.NewServer(&opts)
	if err != nil {
		panic(err)
	}
	go natsd.Start()

	if !natsd.ReadyForConnections(4 * time.Second) {
		panic("not ready for connection")
	}

	nc, err := nats.Connect(natsd.ClientURL())

	if err != nil {
		panic(err)
	}

	subject := "my-subject"

	// TODO: Clean this up when we got this nats stuff working
	nc.Subscribe(subject, func(msg *nats.Msg) {
		data := string(msg.Data)
		fmt.Println("From subscription:", data)
		//ns.Shutdown()
	})

	nc.Publish(subject, []byte("Hello embedded NATS!"))
	//ns.WaitForShutdown()
}
