package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/dronemill/eventsocket"
	"github.com/gorilla/websocket"
)

var d websocket.Dialer
var writer = flag.Bool("writer", false, "Do we write to the websocket connection")

func init() {
	d = websocket.Dialer{
		NetDial:          nil,
		TLSClientConfig:  nil,
		HandshakeTimeout: time.Second * 5,
		ReadBufferSize:   4096,
		WriteBufferSize:  4096,
	}
}

func main() {
	flag.Parse()

	resp, err := http.Post("http://127.0.0.1:8080/v1/clients", "application/json", strings.NewReader(""))
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	client := eventsocket.Client{}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	json.Unmarshal(body, &client)

	headers := http.Header{}
	c, _, err := d.Dial(fmt.Sprintf("ws://127.0.0.1:8080/v1/clients/%s/ws", client.Id), headers)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Connected!  %v  <====>  %v\n\n", c.RemoteAddr(), c.LocalAddr())

	if *writer {
		go becomeWriter(c)
	}

	for {
		messageType, p, err := c.ReadMessage()
		if err != nil {
			panic(err)
		}

		fmt.Printf("%v :: %s\n", messageType, string(p))
	}
}

func becomeWriter(c *websocket.Conn) {
	tickChan := time.NewTicker(time.Second * 1).C

	for {
		select {
		case <-tickChan:
			p := make(map[string]interface{})
			p["value"] = randomString(32)
			m := eventsocket.Message{
				MessageType: eventsocket.MESSAGE_TYPE_BROADCAST,
				Payload:     p,
			}
			if err := c.WriteJSON(m); err != nil {
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
