package eventsocket

import (
	"fmt"

	"github.com/dronemill/eventsocket/Godeps/_workspace/src/github.com/nu7hatch/gouuid"
)

var uuidBuilder = make(chan uuid.UUID)

func init() {
	go populateUuidBuilder()
}

func populateUuidBuilder() {
	for {
		uid, err := uuid.NewV4()
		if err != nil {
			panic(fmt.Sprintf("Error getting an uuid:", err))
		}

		uuidBuilder <- *uid
	}
}
