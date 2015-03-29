package eventsocket

import (
	"testing"

	"github.com/nu7hatch/gouuid"
)

func Test_Uuid_GetUuid(t *testing.T) {

	var uuids map[uuid.UUID]bool
	uuids = make(map[uuid.UUID]bool, 128)

	for i := 0; i < 128; i++ {
		uuid := <-uuidBuilder

		if _, ok := uuids[uuid]; ok {
			t.Fatalf("Duplicate UUID: %s", uuid.String())
		}

		uuids[uuid] = true
	}

}
