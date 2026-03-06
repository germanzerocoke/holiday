package repository

import (
	"backend/online/server/document"
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *Repository) SaveNewMemberId(idRaw []byte) error {
	id := bson.Binary{
		Subtype: 0x04,
		Data:    idRaw,
	}
	doc := document.Member{Id: id}
	result, err := r.db.Collection("member").InsertOne(context.Background(), doc)
	if err != nil {
		slog.Error("fail to save new member id",
			"err", err,
			"id.Data", id.Data,
		)
		return err
	}
	slog.Info("success to save new member id",
		"result", result,
	)

	return nil
}
