package controller

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"backend/online/server/dto"
	"strconv"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func conversationRouter(c *Controller) {
	c.Router(POST, "/conversation/create", c.createConversation)
	c.Router(GET, "/conversation/list?page={page}", c.getConversations)
	c.Router(GET, "/conversation/join?conversationId={conversationId}", c.joinConversation)
}

func (c *Controller) createConversation(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateConversationRequest
	memberIdRaw := r.Header.Get("X-User-Id")
	memberId, err := uuid.Parse(memberIdRaw)
	if err != nil {
		slog.Error("fail to parse userId from X-User-Id header",
			"err", err,
			"memberIdRaw", memberIdRaw,
		)
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, err = w.Write([]byte("incorrect header"))
		if err != nil {
			slog.Error("fail to write response body",
				"err", err,
			)
		}
		return
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Info("incorrect body",
			"err", err,
		)
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, err = w.Write([]byte("incorrect request body"))
		if err != nil {
			slog.Error("fail to write response body",
				"err", err,
			)
		}
		return
	}

	length, err := time.ParseDuration(req.Length)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		slog.Error("fail to parse duration from rawLength",
			"err", err,
			"req.Length", req.Length,
		)
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
		length,
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
	}
}

func (c *Controller) getConversations(w http.ResponseWriter, r *http.Request) {
	memberIdRaw := r.Header.Get("X-User-Id")
	memberId, err := uuid.Parse(memberIdRaw)
	if err != nil {
		slog.Error("fail to parse userId from X-User-Id header",
			"err", err,
			"memberIdRaw", memberIdRaw,
		)
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, err = w.Write([]byte("incorrect header"))
		if err != nil {
			slog.Error("fail to write response body",
				"err", err,
			)
		}
		return
	}
	pageRaw := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageRaw)
	if err != nil {
		slog.Info("incorrect query param for page",
			"err", err,
			"pageRaw", pageRaw)
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, err = w.Write([]byte("wrong cursor"))
		if err != nil {
			slog.Error("fail to write response body",
				"err", err,
			)
		}
		return
	}
	result, err := c.service.GetConversations(r.Context(), memberId, page)
	if err != nil {
		w.WriteHeader(getStatusCode(err))
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			slog.Error("fail to write response body",
				"err", err)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		slog.Error("fail to write response body",
			"err", err)
	}
}

func (c *Controller) joinConversation(w http.ResponseWriter, r *http.Request) {
	memberId := r.Header.Get("X-User-Id")
	conversationIdRaw := r.URL.Query().Get("conversationId")
	conversationId, err := bson.ObjectIDFromHex(conversationIdRaw)
	if err != nil {
		slog.Error("fail to parse conversation object id from raw string",
			"conversationIdRaw", conversationIdRaw)
		return
	}
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		//OriginPatterns:     []string{"example.com"},
		InsecureSkipVerify: true,
	})
	if err != nil {
		slog.Info("fail to accept connection", err)
		return
	}

	ip, err := getPodIP()
	if err != nil {
		slog.Error("fail to get pod ip",
			"err", err)
		return
	}

	defer func() {
		conn.Close(websocket.StatusNormalClosure, "")
		delete(c.room[conversationIdRaw], conn)
		if len(c.room[conversationIdRaw]) == 0 {
			delete(c.room, conversationIdRaw)
			c.service.RemoveServerIP(r.Context(), conversationId, ip)
		}
	}()

	if c.room[conversationIdRaw] == nil {
		c.room[conversationIdRaw] = make(map[*websocket.Conn]struct{})
		err = c.service.AddServerIP(r.Context(), conversationId, ip)
		if err != nil {
			return
		}
	}

	c.room[conversationIdRaw][conn] = struct{}{}

	for {
		ctx := context.Background()
		msgType, data, err := conn.Read(ctx)
		if msgType != websocket.MessageText {
			slog.Error("incorrect message type")
			return
		}
		if websocket.CloseStatus(err) != -1 {
			slog.Error("Connection closed",
				"err", err)
			return
		}
		if err != nil {
			slog.Error("read error",
				"err", err)
			return
		}

		resp := dto.ConversationSignalResponse{
			MemberId: memberId,
			Signal:   data,
		}

		payload, err := json.Marshal(resp)
		if err != nil {
			slog.Error("fail to marshalling conversationSignal",
				"resp", resp)
		}

		for client := range c.room[conversationIdRaw] {
			if client == conn {
				continue
			}
			err = client.Write(ctx, websocket.MessageText, payload)
			if err != nil {
				slog.Error("fail to relay signal")
				return
			}
		}
		err = c.service.PublishConversationSignal(ip, conversationIdRaw, memberId, data)
		if err != nil {
			return
		}
	}
}

func getPodIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String(), nil
		}
	}
	return "", errors.New("IP not found")
}
