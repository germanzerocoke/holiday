package dto

import "encoding/json"

type ConversationSignalMessage struct {
	ServerIP       string          `json:"serverIP"`
	ConversationId string          `json:"conversationId"`
	MemberId       string          `json:"memberId"`
	Signal         json.RawMessage `json:"signal"`
}
