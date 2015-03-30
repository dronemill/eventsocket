package eventsocket

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type controllerClient struct {
	C *httpController
}

// handle a request to create a client
func (C *controllerClient) Create(w http.ResponseWriter, r *http.Request) {
	// close the body fd when we exit
	defer r.Body.Close()

	// create the client
	client := newClient()

	json.NewEncoder(w).Encode(client)

	return
}

// handle a request to create a client
func (C *controllerClient) ServeWs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	client, err := clientById(vars["id"])
	if err != nil {
		http.Error(w, "Client not found", 404)
		return
	}

	client.connectionUpgrade(w, r)

	// json.NewEncoder(w).Encode(client)

	return
}
