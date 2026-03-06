package network

import (
	"backend/auth/config"
	"backend/auth/server/service"

	"github.com/gin-gonic/gin"
)

type Network struct {
	service *service.Service
	engine  *gin.Engine
	port    string
}

func NewNetwork(cfg *config.Config, s *service.Service) *Network {
	n := &Network{
		service: s,
		engine:  gin.New(),
		port:    cfg.Info.Port,
	}

	setGin(n.engine)

	emailRouter(n)
	tokenRouter(n)
	smsRouter(n)

	return n
}

func (n *Network) Start() error {
	return n.engine.Run(n.port)
}
