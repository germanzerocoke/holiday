package controller

import "net/http"

func clubRouter(c *Controller) {
	c.Router(GET, "/club", c.createClub)
}

func (c *Controller) createClub(writer http.ResponseWriter, request *http.Request) {

}
