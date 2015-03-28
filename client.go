package eventsocket

type Client struct {
	Id string `json:"Id"`
}

type Clients []*Client

func newClient() (client *Client) {
	client = new(Client)

	id := <-uuidBuilder
	client.Id = id.String()

	clients = append(clients, client)

	return
}

var clients = make(Clients, 16)
