package repository

import (
	"context"
	"log/slog"
	"online/server/document"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *Repository) SaveConversation(ctx context.Context, memberId uuid.UUID, conversationId bson.ObjectID, novel, shortStory, poem, drama, film, by, rule string, capacity int, when time.Time, length time.Duration) error {
	newConversation := document.Conversation{
		Id:         conversationId,
		Novel:      novel,
		ShortStory: shortStory,
		Poem:       poem,
		Drama:      drama,
		Film:       film,
		By:         by,
		Rule:       rule,
		Capacity:   capacity,
		When:       when,
		Length:     length,
		ModeratorIds: []bson.Binary{
			{4, memberId[:]},
		},
	}
	filter := bson.M{"_id": bson.Binary{Subtype: 4, Data: memberId[:]}}
	update := bson.M{
		"$push": bson.M{
			"m_c_ids": conversationId,
		},
	}

	session, err := r.client.StartSession()
	if err != nil {
		slog.Error("fail to start transaction for insert new conversation",
			"err", err)
		return err
	}
	defer session.EndSession(ctx)
	_, err = session.WithTransaction(ctx, func(ctx context.Context) (any, error) {
		_, err = r.db.Collection("conversation").InsertOne(ctx, newConversation)
		if err != nil {
			slog.Error("fail to insert new conversation",
				"err", err)
			return nil, err
		}

		_, err = r.db.Collection("member").UpdateOne(ctx, filter, update)
		if err != nil {
			slog.Error("fail to insert new conversation id to moderator conversation member array",
				"err", err)
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		slog.Error("fail to transaction for saving new conversation",
			"err", err,
		)
		return err
	}

	return nil
}
