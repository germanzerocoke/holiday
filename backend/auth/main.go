package main

import (
	"backend/auth/config"
	"backend/auth/server"
	"flag"
)

var cfgPath = flag.String("cfg", "./config.toml", "config path")

func main() {
	flag.Parse()

	cfg := config.NewConfig(*cfgPath)

	server.NewServer(cfg)
}
