package server

import (
	"net/http"
	"online/server/controller"
	"online/server/kafka/consumer"
	"online/server/kafka/producer"
	"online/server/repository"
	"online/server/service"
)

func NewServer(mux *http.ServeMux) {

	kp := producer.NewKafkaProducer()

	r := repository.NewRepository()

	s := service.NewService(r, kp)

	ks := consumer.NewKafkaConsumer(s)

	go func() {
		ks.GetMessage([]string{"auth.new_member_id"})
	}()

	controller.SetController(s, mux)
}
