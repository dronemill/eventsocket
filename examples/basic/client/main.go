package main

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"time"

	"github.com/dronemill/eventsocket-client-go"
)

var writer = flag.Bool("writer", false, "Do we write to the websocket connection")

func main() {
	flag.Parse()

	client, err := eventsocketclient.NewClient("127.0.0.1:8080")
	if err != nil {
		panic(err)
	}

	if err := client.DialWs(); err != nil {
		panic(err)
	}

	fmt.Println("Connected!")

	if *writer {
		go becomeWriter(client)
	}

	client.Suscribe("foo")

	for {
		messageType, p, err := client.ReadMessage()
		if err != nil {
			panic(err)
		}

		fmt.Printf("%v :: %s", messageType, string(p))
	}
}

func becomeWriter(client *eventsocketclient.Client) {
	tickChan := time.NewTicker(time.Second).C

	for {
		select {
		case <-tickChan:
			p := eventsocketclient.NewPayload()
			p["awesomeValue"] = randomString(32)

			if err := client.Emit("foo", &p); err != nil {
				panic(err)
			}

			// p := make(eventsocketclient.Payload)
			// p["value"] = randomString(32)

			// if err := client.Broadcast(&p); err != nil {
			// 	panic(err)
			// }
		}
	}
}

func randomString(size int) string {
	rb := make([]byte, size)
	_, err := rand.Read(rb)

	if err != nil {
		fmt.Println(err)
	}

	return base64.URLEncoding.EncodeToString(rb)
}
