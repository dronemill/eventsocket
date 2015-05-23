package eventsocket

import (
	"errors"
	"fmt"
	"net/http"

	log "github.com/dronemill/eventsocket/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/dronemill/eventsocket/Godeps/_workspace/src/github.com/gorilla/websocket"
)

type Client struct {
	Id            string          `json:"Id"`
	ws            *wsConnection   `json:-`
	subscriptions map[string]bool `json:-`
}

type Clients map[string]*Client

// the main client store
var clients = make(Clients)

// instantiate a new client, set it's id, and store the client
func newClient() (client *Client) {

	client = new(Client)

	id := <-uuidBuilder
	client.Id = id.String()
	client.subscriptions = make(map[string]bool)

	clients[client.Id] = client

	log.WithField("clientID", id).Info("Created new Client")
	return
}

// Connection upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// FIXME: this is bad
	CheckOrigin: checkOrigin,
}

func checkOrigin(r *http.Request) bool {
	return true
}

// fetch a client by it's id
func clientById(id string) (*Client, error) {
	if client, ok := clients[id]; ok {
		return client, nil
	}
	return nil, errors.New(fmt.Sprintf("Client id does not exist: %s", id))
}

// upgrade the http connection to become a ws connection
func (client *Client) connectionUpgrade(w http.ResponseWriter, r *http.Request) error {
	// sanity check
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return errors.New("Methods not allowed")
	}

	// upgrade the connection
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	ws, err := newWsConnection(wsConn)
	if err != nil {
		return err
	}

	// store the connection reference
	client.ws = ws

	h.register <- client

	return nil
}

// receive message on behalf of the client
func (client *Client) recv() {
	for {
		// get a message from the channel
		message := <-client.ws.recv

		// if the ws was closed, then get our
		if client.ws.closed {
			return
		}

		h.recvClientMessage <- &ClientMessage{
			ClientId: client.Id,
			Message:  message,
		}
	}

}

func (client *Client) run() error {

	// start the client message handler
	go client.recv()

	// socket read write control
	client.ws.pump()

	return nil
}
