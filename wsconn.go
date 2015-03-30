package eventsocket

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type wsConnection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

func newWsConnection(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return errors.New("Methods not allowed")
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	c := &wsConnection{send: make(chan []byte, 256), ws: ws}

	fmt.Printf("%+v\n", c)

	// h.register <- c
	// go c.writePump()
	// c.readPump()

	return nil
}
