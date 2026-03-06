package server

import (
	"caller/server/kafka/consumer"
	"caller/server/repository"
	"caller/server/service"
)

func NewServer() {
	r := repository.NewRepository()

	s := service.NewService(r)

	ks := consumer.NewKafkaConsumer(s)

	ks.GetMessage([]string{"conversation.signal"})

}
