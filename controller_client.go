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

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

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

	err = client.connectionUpgrade(w, r)
	if err != nil {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	client.ws.pump()

	// json.NewEncoder(w).Encode(client)

	return
}
