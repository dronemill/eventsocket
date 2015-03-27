package eventsocket

type httpController struct {
	Client *controllerClient
}

func newHttpController() (c *httpController, err error) {
	c = &httpController{
		Client: &controllerClient{c},
	}

	return
}
