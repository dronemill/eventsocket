package eventsocket

import (
	"errors"
	"fmt"
	"net/http"
)

type Client struct {
	Id string        `json:"Id"`
	ws *wsConnection `json:-`
}

type Clients map[string]*Client

// the main client store
var clients = make(Clients)

// instantiate a new client, set it's id, and store the client
func newClient() (client *Client) {
	client = new(Client)

	id := <-uuidBuilder
	client.Id = id.String()

	clients[client.Id] = client

	return
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
	return nil
}
