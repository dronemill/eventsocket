package eventsocket

import (
	"encoding/json"
	"net/http"
)

type controllerClient struct {
	C *httpController
}

// handle a request to create a client
func (C *controllerClient) Create(w http.ResponseWriter, r *http.Request) {
	// log.Info("Creating a client")

	// close the body fd when we exit
	defer r.Body.Close()

	// create the client
	client := newClient()

	json.NewEncoder(w).Encode(client)

	return
}
