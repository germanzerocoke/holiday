package document

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Club struct {
	Id          bson.ObjectID `bson:"_id"`
	Name        string        `bson:"name"`
	Description string        `bson:"description"`

	//Members are club independent so reference
	//since we use uuid in token claim subject we use bson.Binary instead of bson.ObjectID
	ModIds        []bson.Binary `bson:"mod_ids"`
	SubscriberIds []bson.Binary `bson:"subscriber_ids"`
	//this is for blocking participation at the further meetings,
	//listening(subscribing) can't be blocked
	BlockedIds []bson.Binary `bson:"blocked_ids"`
	//QNAs are club subset so it seemed good to embedding,
	//but embedding make query writing difficult, so we will ref it
	MeetingIds []bson.ObjectID `bson:"meeting_ids"`
}

type Meeting struct {
	//topic and id will send to elastic search
	Id       bson.ObjectID `bson:"_id"`
	Topic    string        `bson:"topic"`
	Ground   string        `bson:"ground"`
	rule     string        `bson:"rule"`
	Capacity int           `bson:"capacity"`
	When     time.Time     `bson:"when"`
	Length   time.Duration `bson:"length"`
	Done     bool          `bson:"done"`

	ClubId         bson.ObjectID   `bson:"club_id"`
	ModId          bson.Binary     `bson:"mod_id"`
	ParticipantIds []bson.Binary   `bson:"participant_ids"`
	QNAIds         []bson.ObjectID `bson:"qna_ids"`

	//AudioFileUrl string        `bson:"audio_file_url"`
}

type QNA struct {
	Id       bson.ObjectID `bson:"_id"`
	Question string        `bson:"question"`

	MeetingId bson.ObjectID   `bson:"meeting_id"`
	AnswerIds []bson.ObjectID `bson:"answer_ids"`
}

type Answer struct {
	Id      bson.ObjectID `bson:"_id"`
	Content string        `bson:"content"`
}
