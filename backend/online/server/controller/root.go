package controller

import (
	"net/http"
	"online/server/service"

	"github.com/coder/websocket"
)

type HTTPMethod int

const (
	GET HTTPMethod = iota
	POST
	DELETE
	PUT
)

type Controller struct {
	service *service.Service
	mux     *http.ServeMux
	room    map[string]map[*websocket.Conn]struct{}
}

func SetController(s *service.Service, m *http.ServeMux) {

	c := &Controller{
		service: s,
		mux:     m,
	}

	conversationRouter(c)
}

func (c *Controller) Router(httpMethod HTTPMethod, path string, handler http.HandlerFunc) {
	m := c.mux

	switch httpMethod {
	case GET:
		m.HandleFunc("GET "+path, handler)
	case POST:
		m.HandleFunc("POST "+path, handler)
	case PUT:
		m.HandleFunc("PUT "+path, handler)
	case DELETE:
		m.HandleFunc("DELETE "+path, handler)

	default:
		panic("This HTTP method is not supported")
	}
}

func getStatusCode(err error) int {
	return http.StatusBadRequest
}
