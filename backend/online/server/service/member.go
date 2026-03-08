package service

import (
	"context"

	"github.com/google/uuid"
)

func (s *Service) SaveNewMemberId(idRaw []byte) error {
	err := s.repository.SaveNewMemberId(idRaw)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) SetServerIP(ctx context.Context, memberId uuid.UUID, ip string) error {
	err := s.repository.SetServerIP(ctx, memberId, ip)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) RemoveServerIP(ctx context.Context, memberId uuid.UUID) error {
	err := s.repository.RemoveServerIP(ctx, memberId)
	if err != nil {
		return err
	}
	return nil
}
