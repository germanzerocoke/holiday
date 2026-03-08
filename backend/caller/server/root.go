package server

import (
	"backend/caller/server/kafka/consumer"
	"backend/caller/server/repository"
	"backend/caller/server/service"
	"log/slog"
	"os"
)

func NewServer() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	r := repository.NewRepository()

	s := service.NewService(r)

	ks := consumer.NewKafkaConsumer(s)

	ks.GetMessage([]string{"conversation.signal"})
}
