package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/dronemill/eventsocket"
	"github.com/gorilla/websocket"
)

var d websocket.Dialer

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

	fmt.Printf("Connected!  %v  <====>  %v", c.RemoteAddr(), c.LocalAddr())
}
