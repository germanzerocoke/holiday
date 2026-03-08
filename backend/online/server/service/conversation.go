package service

import (
	"backend/online/server/dto"
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *Service) CreateConversation(
	ctx context.Context,
	memberId uuid.UUID,
	novel,
	shortStory,
	poem,
	drama,
	film,
	by,
	rule string,
	capacity int,
	when time.Time,
	length time.Duration,
) (map[string]string, error) {
	conversationId := bson.NewObjectID()
	err := s.repository.SaveConversation(
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

func (s *Service) GetConversations(ctx context.Context, memberId uuid.UUID, page int) (resp []dto.ConversationFeedResponse, err error) {
	items, err := s.repository.GetNextConversations(ctx, page)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		var isModerator bool
		for _, mId := range item.ModeratorIds {
			if bytes.Equal(mId.Data, memberId[:]) {
				isModerator = true
			}
		}
		var isRegistrant bool
		for _, pId := range item.RegistrantIds {
			if bytes.Equal(pId.Data, memberId[:]) {
				isRegistrant = true
			}
		}
		var onAir bool
		if time.Now().After(item.When.Add(-10 * time.Minute)) {
			onAir = true
		}

		resp = append(resp, dto.ConversationFeedResponse{
			Id:           item.Id.Hex(),
			Novel:        item.Novel,
			ShortStory:   item.ShortStory,
			Poem:         item.Poem,
			Drama:        item.Drama,
			Film:         item.Film,
			By:           item.By,
			Rule:         item.Rule,
			When:         item.When,
			Length:       item.Length.String(),
			OnAir:        onAir,
			IsModerator:  isModerator,
			IsRegistrant: isRegistrant,
		})
	}
	return resp, nil
}

func (s *Service) GetParticipants(ctx context.Context, conversationId bson.ObjectID) (pids []string, err error) {
	pidsRaw, err := s.repository.GetParticipants(ctx, conversationId)
	if err != nil {
		return nil, err
	}
	for _, pidRaw := range pidsRaw {
		pid, err := uuid.FromBytes(pidRaw.Data)
		if err != nil {
			slog.Error("fail to parse uuid from pidRaw",
				"err", err,
				"pidRaw.Data", pidRaw.Data)
			return nil, err
		}
		pids = append(pids, pid.String())
	}
	return pids, nil
}

func (s *Service) AddParticipant(ctx context.Context, conversationId bson.ObjectID, memberId uuid.UUID) error {
	err := s.repository.AddParticipant(ctx, conversationId, memberId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) RemoveParticipant(ctx context.Context, conversationId bson.ObjectID, memberId uuid.UUID) error {
	err := s.repository.RemoveParticipant(ctx, conversationId, memberId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) PublishConversationSignal(fromId, toId string, data []byte) error {
	msg := dto.ConversationSignalMessage{
		FromId: fromId,
		ToId:   toId,
		Signal: data,
	}

	value, err := json.Marshal(msg)
	if err != nil {
		slog.Error("fail to marshal msg",
			"err", err,
			"msg", msg)
		return err
	}
	err = s.kafkaProducer.PushMessage("conversation.signal", value)
	if err != nil {
		return err
	}

	return nil
}
