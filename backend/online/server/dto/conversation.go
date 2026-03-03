package dto

import (
	"time"
)

type CreateConversationRequest struct {
	Novel      string    `json:"novel"`
	ShortStory string    `json:"short_story"`
	Poem       string    `json:"poem"`
	Drama      string    `json:"drama"`
	Film       string    `json:"film"`
	By         string    `json:"by"`
	Rule       string    `json:"rule"`
	Capacity   int       `json:"capacity"`
	When       time.Time `json:"when"`
	Length     string    `json:"length"`
}

type ConversationFeedResponse struct {
	Id           string    `json:"id"`
	Novel        string    `json:"novel,omitempty"`
	ShortStory   string    `json:"shortStory,omitempty"`
	Poem         string    `json:"poem,omitempty"`
	Drama        string    `json:"drama,omitempty"`
	Film         string    `json:"film,omitempty"`
	By           string    `json:"by,omitempty"`
	Rule         string    `json:"rule,omitempty"`
	When         time.Time `json:"when"`
	Length       string    `json:"length"`
	OnAir        bool      `json:"onAir"`
	IsModerator  bool      `json:"isModerator"`
	IsRegistrant bool      `json:"isRegistrant"`
}
