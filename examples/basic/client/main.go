package main

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"time"

	"github.com/dronemill/eventsocket-client-go"
)

var ticker = flag.Bool("ticker", false, "Do we write to the websocket connection as we tick?")
var requester = flag.Bool("requester", false, "Do we write to the websocket connection?")
var requestClientId = flag.String("requestClientId", "", "The clientId to send a request to")

func main() {
	flag.Parse()

	client, err := eventsocketclient.NewClient("127.0.0.1:8080")
	if err != nil {
		panic(err)
	}

	if err := client.DialWs(); err != nil {
		panic(err)
	}

	fmt.Printf("Connected! ClientId: %s\n", client.Id)

	if *ticker {
		go tick(client)
	}

	if *requester {
		go makeRequest(client)
	}

	go client.Recv()

	subChan, err := client.Suscribe("foo")
	if err != nil {
		panic(err)
	}

	for {
		select {
		case m := <-client.RecvBroadcast:
			if m.Err != nil {
				panic(m.Err)
			}
			fmt.Printf("BROADCAST err:%v :: %+v\n", m.Err, m.Message.Payload)
		case m := <-subChan:
			if m.Err != nil {
				panic(m.Err)
			}
			fmt.Printf("Foo: err:%v :: %+v\n", m.Err, m.Message.Payload)
		case m := <-client.RecvRequest:
			if m.Err != nil {
				panic(m.Err)
			}
			fmt.Printf("REQUEST err:%v :: %+v\n", m.Err, m.Message.Payload)
			p := eventsocketclient.NewPayload()
			p["I_SEE_YOU"] = randomString(32)
			if err := client.Reply(m.Message.RequestId, m.Message.ReplyClientId, &p); err != nil {
				panic(err)
			}
		}
	}
}

func makeRequest(client *eventsocketclient.Client) {
	p := eventsocketclient.NewPayload()
	p["PEEK_ABOO"] = randomString(32)

	res, err := client.Request(*requestClientId, &p)
	if err != nil {
		panic(err)
	}

	m := <-res

	fmt.Printf("REPLY err:%v :: %+v\n", m.Err, m.Message.Payload)
}

func tick(client *eventsocketclient.Client) {
	tickChan := time.NewTicker(time.Second).C

	for {
		select {
		case <-tickChan:
			p := eventsocketclient.NewPayload()
			p["awesomeValue"] = randomString(32)

			if err := client.Emit("foo", &p); err != nil {
				panic(err)
			}

			p = make(eventsocketclient.Payload)
			p["value"] = randomString(32)

			if err := client.Broadcast(&p); err != nil {
				panic(err)
			}
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
