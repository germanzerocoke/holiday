package service

import "online/server/repository"

type Service struct {
	repository *repository.Repository
}

func NewService(r *repository.Repository) *Service {
	s := &Service{repository: r}
	return s
}
