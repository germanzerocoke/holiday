package server

import (
	"backend/caller/server/kafka/consumer"
	"backend/caller/server/repository"
	"backend/caller/server/service"
)

func NewServer() {
	r := repository.NewRepository()

	s := service.NewService(r)

	ks := consumer.NewKafkaConsumer(s)

	ks.GetMessage([]string{"conversation.signal"})
}
