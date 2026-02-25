package controller

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"online/server/dto"
)

func clubRouter(c *Controller) {
	c.Router(GET, "/club", c.createClub)
}

func (c *Controller) createClub(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateClubRequest
	userId := r.Header.Get("X-User-Id")
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("fail to close request body",
				"err", err,
			)
		}
	}(r.Body)
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
	result, err := c.service.CreateClub(r.Context(), userId, req.Name, req.Description)
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
