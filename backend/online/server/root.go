package server

import (
	"backend/online/server/controller"
	"backend/online/server/kafka/consumer"
	"backend/online/server/kafka/producer"
	"backend/online/server/repository"
	"backend/online/server/service"
	"net/http"
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
