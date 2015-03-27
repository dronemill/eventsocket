package eventsocket

import "net/http"

type controllerClient struct {
	C *httpController
}

func (C *controllerClient) Create(w http.ResponseWriter, r *http.Request) {
	// log.Info("Creating a client")

	// close the body fd when we exit
	defer r.Body.Close()

	panic("Awwww yeahhhh")
}
