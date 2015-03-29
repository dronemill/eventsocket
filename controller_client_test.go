package eventsocket

import (
	"net/http"
	"strings"
	"testing"
)

func Test_ControllerClient_Create(t *testing.T) {
	server, w := testHttpServer(t)

	req, _ := http.NewRequest("POST", "/v1/clients", strings.NewReader(""))
	server.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Client create didn't return %v. Got: %v", http.StatusOK, w.Code)
	}

	// TODO do something more than just check status code
	return
}
