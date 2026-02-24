package controller

import (
	"net/http"
	"online/server/service"
)

type Controller struct {
	service *service.Service
	mux     *http.ServeMux
}

func NewController(s *service.Service, m *http.ServeMux) *Controller {

	c := &Controller{
		service: s,
	}
	return c
}
