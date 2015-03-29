package eventsocket

import (
	"net/http"
	"testing"
)

func Test_ControllerDev_Ping(t *testing.T) {
	server, w := testHttpServer(t)

	req, _ := http.NewRequest("GET", "/v1/dev/ping", nil)
	server.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected %v. Got: %v", http.StatusOK, w.Code)
	}
}
