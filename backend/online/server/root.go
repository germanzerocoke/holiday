package server

import (
	"net/http"
	"online/server/controller"
	"online/server/kafka/consumer"
	"online/server/repository"
	"online/server/service"
)

func NewServer(mux *http.ServeMux) {
	r := repository.NewRepository()

	s := service.NewService(r)

	k := consumer.NewKafkaConsumer(s)

	go func() {
		k.GetMessage([]string{"auth.new_member_id"})
	}()

	controller.SetController(s, mux)
}
