package eventsocket

import (
	"reflect"
	"testing"
)

func Test_Client_newClient(t *testing.T) {
	client := newClient()

	if reflect.TypeOf(client).String() != "*eventsocket.Client" {
		t.Fatal("Client is not of type *eventsocket.Client")
	}
}
