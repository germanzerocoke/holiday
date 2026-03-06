package document

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Conversation struct {
	Id         bson.ObjectID `bson:"_id"`
	Novel      string        `bson:"novel,omitempty"`
	ShortStory string        `bson:"short_story,omitempty"`
	Poem       string        `bson:"poem,omitempty"`
	Drama      string        `bson:"drama,omitempty"`
	Film       string        `bson:"film,omitempty"`
	By         string        `bson:"by,omitempty"`
	Rule       string        `bson:"rule,omitempty"`
	Capacity   int           `bson:"capacity"`
	When       time.Time     `bson:"when"`
	Length     time.Duration `bson:"length"`
	Expired    bool          `bson:"expired"`

	ModeratorIds  []bson.Binary `bson:"m_ids"`
	RegistrantIds []bson.Binary `bson:"r_ids"`
	ServerIPs     []string      `bson:"s_ips"`
}
