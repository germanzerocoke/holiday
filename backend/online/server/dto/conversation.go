package dto

import "time"

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
