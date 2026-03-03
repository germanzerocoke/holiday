package service

import (
	"bytes"
	"context"
	"online/server/dto"
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
