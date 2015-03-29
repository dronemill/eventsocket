package eventsocket

import "net/http"

type controllerDev struct {
	C *httpController
}

func (C *controllerDev) Ping(w http.ResponseWriter, r *http.Request) {
	// log.Info("Creating a client")

	// close the body fd when we exit
	defer r.Body.Close()

	w.Write([]byte("pong"))
}
