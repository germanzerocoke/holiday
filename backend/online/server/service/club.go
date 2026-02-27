package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/IBM/sarama"
)

func (s *Service) CreateClub(
	ctx context.Context, userId, name, description string) (
	map[string]string, error) {

}
