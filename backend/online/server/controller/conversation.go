package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"online/server/dto"
)

func conversationRouter(c *Controller) {
	c.Router(GET, "/conversation/create", c.createConversation)
}

func (c *Controller) createConversation(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateConversationRequest

	memberId := r.Header.Get("X-User-Id")
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Info("incorrect body",
			"err", err,
		)
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, err = w.Write([]byte("wrong input"))
		if err != nil {
			slog.Error("fail to write response body",
				"err", err,
			)
		}
		return
	}
	result, err := c.service.CreateConversation(
		r.Context(),
		memberId,
		req.Novel,
		req.ShortStory,
		req.Poem,
		req.Drama,
		req.Film,
		req.By,
		req.Rule,
		req.Capacity,
		req.When,
		req.Length,
	)
	if err != nil {
		w.WriteHeader(getStatusCode(err))
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			slog.Error("fail to write response body",
				"err", err,
			)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		slog.Error("fail to write response body",
			"err", err,
		)
		return
	}
	return
}
