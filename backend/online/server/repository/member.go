package repository

import (
	"backend/online/server/document"
	"context"
	"log/slog"

	"github.com/google/uuid"
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

func (r *Repository) SetServerIP(ctx context.Context, memberId uuid.UUID, ip string) error {
	_, err := r.db.Collection("member").
		UpdateOne(ctx, bson.M{"_id": bson.Binary{Subtype: 4, Data: memberId[:]}},
			bson.M{"$set": bson.M{"server_ip": ip}})
	if err != nil {
		slog.Error("fail to set ip to conversation doc's server ips",
			"err", err)
		return err
	}
	return nil
}

func (r *Repository) RemoveServerIP(ctx context.Context, memberId uuid.UUID) error {

	_, err := r.db.Collection("member").
		UpdateOne(ctx, bson.M{"_id": bson.Binary{Subtype: 4, Data: memberId[:]}},
			bson.M{"$unset": bson.M{"server_ip": ""}})
	if err != nil {
		slog.Error("fail to remove ip to conversation doc's server ips",
			"err", err)
		return err
	}
	return nil
}
