package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *Service) CreateConversation(
	ctx context.Context,
	memberIdRaw,
	novel,
	shortStory,
	poem,
	drama,
	film,
	by,
	rule string,
	capacity int,
	when time.Time,
	lengthRaw string,
) (map[string]string, error) {
	memberId, err := uuid.Parse(memberIdRaw)
	if err != nil {
		slog.Error("fail to parse userId from X-User-Id header",
			"err", err,
			"memberIdRaw", memberIdRaw,
		)
		return nil, err
	}
	length, err := time.ParseDuration(lengthRaw)
	if err != nil {
		slog.Error("fail to parse duration from rawLength",
			"err", err,
			"rawLength", lengthRaw,
		)
		return nil, err
	}
	conversationId := bson.NewObjectID()
	err = s.repository.SaveConversation(
		ctx,
		memberId,
		conversationId,
		novel,
		shortStory,
		poem,
		drama,
		film,
		by,
		rule,
		capacity,
		when,
		length,
	)
	if err != nil {
		return nil, err
	}

	return map[string]string{"conversationId": conversationId.Hex()}, nil

}
