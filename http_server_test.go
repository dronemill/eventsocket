package eventsocket

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_HttpServer_Router(t *testing.T) {
	server := &httpServer{}

	err := server.route()
	if err != nil {
		t.Fatalf("Received error while installing router %s", err.Error())
	}

	req, _ := http.NewRequest("GET", "/v1/dev/ping", nil)

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Dev ping page didn't return %v. Got: %v", http.StatusOK, w.Code)
	}

	if s := w.Body.String(); s != "pong" {
		t.Fatalf("Expected \"pong\" got \"%s\"", s)
	}
}
