package eventsocket

import (
	"time"

	"github.com/dronemill/eventsocket/Godeps/_workspace/src/github.com/gorilla/websocket"
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

type wsConnection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan Message

	// Receive messages from this client
	recv chan Message

	// has this wsConnection been closed
	closed bool
}

func newWsConnection(ws *websocket.Conn) (*wsConnection, error) {
	wsc := &wsConnection{
		send:   make(chan Message, 256),
		recv:   make(chan Message, 256),
		ws:     ws,
		closed: false,
	}

	return wsc, nil
}

func (wsc *wsConnection) pump() {
	go wsc.writePump()
	wsc.readPump()
}

// readPump pumps messages from the websocket connection to the hub.
func (wsc *wsConnection) readPump() {
	defer func() {
		h.unregister <- wsc
	}()

	wsc.ws.SetReadLimit(maxMessageSize)
	wsc.ws.SetReadDeadline(time.Now().Add(pongWait))
	wsc.ws.SetPongHandler(func(string) error { wsc.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		m := &Message{}
		if wsc.ws.ReadJSON(m) != nil {
			break
		}

		wsc.recv <- *m
	}
}

// write writes a message with the given message type and payload.
func (wsc *wsConnection) write(mt int, payload []byte) error {
	wsc.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return wsc.ws.WriteMessage(mt, payload)
}

// writeJSON writes a json message
func (wsc *wsConnection) writeJSON(message interface{}) error {
	wsc.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return wsc.ws.WriteJSON(message)
}

// writePump pumps messages from the hub to the websocket connection.
func (wsc *wsConnection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		h.unregister <- wsc
	}()
	for {
		select {
		case message, ok := <-wsc.send:
			if !ok {
				wsc.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := wsc.writeJSON(message); err != nil {
				return
			}
		case <-ticker.C:
			if err := wsc.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (wsc *wsConnection) Close() {
	// if we have already closed the channels, then get out
	if wsc.closed {
		return
	}
	wsc.closed = true

	// blose stuff
	close(wsc.send)
	close(wsc.recv)
	wsc.ws.Close()
}
