package server

import (
	"backend/auth/config"
	"backend/auth/server/kafka/producer"
	"backend/auth/server/logger"
	"backend/auth/server/network"
	"backend/auth/server/repository"
	"backend/auth/server/service"
)

type Server struct {
	config *config.Config
}

func NewServer(cfg *config.Config) *Server {
	server := &Server{cfg}

	kp := producer.NewKafkaProducer()

	logger.SetLogger(kp)

	r := repository.NewRepository(cfg)

	s := service.NewService(cfg, r, kp)

	n := network.NewNetwork(cfg, s)

	n.Start()

	return server
}
