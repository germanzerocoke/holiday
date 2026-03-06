package repository

import (
	"caller/server/document"
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (r *Repository) GetServerIPs(ctx context.Context, conversationId bson.ObjectID) ([]string, error) {
	filter := bson.M{
		"_id": conversationId,
	}

	opts := options.FindOne().SetProjection(bson.M{
		"s_ips": 1,
		"_id":   0,
	})

	var d document.Conversation
	err := r.db.Collection("conversation").FindOne(ctx, filter, opts).Decode(&d)
	if err != nil {
		slog.Info("fail to find ips",
			"err", err)
		return nil, err
	}

	return d.ServerIPs, nil
}
