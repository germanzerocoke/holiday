package controller

import "net/http"

func clubRouter(c *Controller) {
	c.Router(GET, "/club", c.createClub)
}

func (c *Controller) createClub(w http.ResponseWriter, r *http.Request) {

}
