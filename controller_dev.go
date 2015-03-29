package eventsocket

import "net/http"

type controllerDev struct {
	C *httpController
}

func (C *controllerDev) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
