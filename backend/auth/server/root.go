package server

import (
	"backend/auth/server/kafka/producer"
	"backend/auth/server/logger"
	"backend/auth/server/network"
	"backend/auth/server/repository"
	"backend/auth/server/service"
)

func NewServer() {

	kp := producer.NewKafkaProducer()

	logger.SetLogger(kp)

	r := repository.NewRepository()

	s := service.NewService(r, kp)

	n := network.NewNetwork(s)

	n.Start()
}
