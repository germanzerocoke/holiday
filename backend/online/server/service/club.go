package service

import (
	"context"

	"github.com/google/uuid"
)

func (s *Service) CreateClub(
	ctx context.Context, userId, name, description string) (
	map[string]string, error) {
	id, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

}
