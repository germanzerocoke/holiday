package server

import (
	"net/http"
	"online/server/controller"
	"online/server/repository"
	"online/server/service"
)

func NewServer(mux *http.ServeMux) {
	r := repository.NewRepository()

	s := service.NewService(r)

	controller.SetController(s, mux)
}
