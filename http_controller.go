package eventsocket

type httpController struct {
	Client *controllerClient
	Dev    *controllerDev
}

func newHttpController() (c *httpController, err error) {
	c = &httpController{
		Client: &controllerClient{c},
		Dev:    &controllerDev{c},
	}

	return
}
