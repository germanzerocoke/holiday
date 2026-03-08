package service

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"
)

func (s *Service) PropagateSignal(ctx context.Context, fromId, toIdRaw string, signal json.RawMessage) error {
	toId, err := uuid.Parse(toIdRaw)
	if err != nil {
		slog.Error("fail to parse")
		return err
	}
	ip, err := s.repository.GetServerIp(ctx, toId)
	if err != nil {
		return err
	}

	return nil
}
