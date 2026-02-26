package document

import "go.mongodb.org/mongo-driver/v2/bson"

type Member struct {
	Id   bson.Binary `bson:"_id"`
	Name string      `bson:"name"`
}
