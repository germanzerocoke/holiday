package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *Service) PropagateSignal(ctx context.Context, ip, conversationIdRaw, memberId string, signal json.RawMessage) error {
	conversationId, err := bson.ObjectIDFromHex(conversationIdRaw)
	if err != nil {
		slog.Error("fail to parse")
		return err
	}
	ips, err := s.repository.GetServerIPs(ctx, conversationId)
	if err != nil {
		return err
	}
	for _, sip := range ips {

	}

}
