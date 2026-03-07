package controller

import (
	"backend/online/server/dto"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net"
	"net/http"
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
	memberIdRaw := r.Header.Get("X-User-Id")
	memberId, err := uuid.Parse(memberIdRaw)
	if err != nil {
		slog.Error("fail to parse member id from raw string",
			"err", err,
			"memberIdRaw", memberIdRaw)
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, err = w.Write([]byte("incorrect header"))
		if err != nil {
			slog.Error("fail to write response body",
				"err", err,
			)
		}
		return
	}
	conversationIdRaw := r.URL.Query().Get("conversationId")
	conversationId, err := bson.ObjectIDFromHex(conversationIdRaw)
	if err != nil {
		slog.Error("fail to parse conversation object id from raw string",
			"conversationIdRaw", conversationIdRaw)
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, err = w.Write([]byte("incorrect header"))
		if err != nil {
			slog.Error("fail to write response body",
				"err", err,
			)
		}
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		//OriginPatterns:     []string{"example.com"},
		InsecureSkipVerify: true,
	})
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, err = w.Write([]byte("incorrect header"))
		if err != nil {
			slog.Error("fail to write response body",
				"err", err,
			)
		}
		slog.Info("fail to accept connection", err)
		return
	}
	c.connections[memberIdRaw] = conn

	ip, _ := getPodIP()

	defer func() {
		destroy := context.Background()
		conn.Close(websocket.StatusNormalClosure, "")
		delete(c.connections, memberIdRaw)
		c.service.RemoveServerIP(destroy, memberId)
		c.service.RemoveParticipant(destroy, conversationId, memberId)
	}()

	init := context.Background()

	err = c.service.SetServerIP(init, memberId, ip)
	if err != nil {
		err = conn.Write(init, websocket.MessageText, []byte(err.Error()))
		if err != nil {
			slog.Error("fail to write payload",
				"err", err,
			)
		}
		return
	}
	pids, err := c.service.GetParticipants(init, conversationId)
	if err != nil {
		err = conn.Write(init, websocket.MessageText, []byte(err.Error()))
		if err != nil {
			slog.Error("fail to write payload",
				"err", err,
			)
		}
		return
	}
	resp := dto.ConversationSignalResponse{ParticipantIds: pids}
	payload, err := json.Marshal(resp)
	if err != nil {
		slog.Error("fail to marshal")
		err = conn.Write(init, websocket.MessageText, []byte("fail to get participants"))
		if err != nil {
			slog.Error("fail to write payload",
				"err", err,
			)
		}
		return
	}
	err = conn.Write(init, websocket.MessageText, payload)
	if err != nil {
		slog.Error("fail to write payload",
			"err", err,
		)
		return
	}

	err = c.service.AddParticipant(init, conversationId, memberId)
	if err != nil {
		err = conn.Write(init, websocket.MessageText, []byte(err.Error()))
		if err != nil {
			slog.Error("fail to write payload",
				"err", err,
			)
		}
		return
	}
	for _, pid := range pids {
		resp = dto.ConversationSignalResponse{
			FromId: memberIdRaw,
		}
		payload, err = json.Marshal(resp)
		if err != nil {
			slog.Error("fail to marshal")
			err = conn.Write(init, websocket.MessageText, []byte("fail to get participants"))
			if err != nil {
				slog.Error("fail to write payload",
					"err", err,
				)
			}
			return
		}
		p, ok := c.connections[pid]
		if ok {
			err = p.Write(init, websocket.MessageText, payload)
			if err != nil {
				slog.Error("fail to write payload",
					"err", err,
				)
				return
			}
			continue
		}
		err = c.service.PublishConversationSignal(memberIdRaw, pid, []byte{})
		if err != nil {
			err = conn.Write(init, websocket.MessageText, []byte("fail to publish"))
			if err != nil {
				slog.Error("fail to write payload",
					"err", err,
				)
			}
			return
		}
	}
	for {
		ctx := context.Background()
		msgType, data, err := conn.Read(ctx)
		if msgType != websocket.MessageText {
			slog.Error("incorrect message type")
			err = conn.Write(ctx, websocket.MessageText, []byte("incorrect message type"))
			if err != nil {
				slog.Error("fail to write payload",
					"err", err,
				)
			}
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
			err = conn.Write(ctx, websocket.MessageText, []byte("read error"))
			if err != nil {
				slog.Error("fail to write payload",
					"err", err,
				)
			}
			return
		}
		var req dto.ConversationSignalRequest
		err = json.Unmarshal(data, &req)
		if err != nil {
			slog.Error("fail to unmarshalling data",
				"err", err)
			err = conn.Write(ctx, websocket.MessageText, []byte("fail to unmarshal"))
			if err != nil {
				slog.Error("fail to write payload",
					"err", err,
				)
			}
			return
		}

		to, ok := c.connections[req.ToId]
		if ok {
			resp = dto.ConversationSignalResponse{
				FromId: memberIdRaw,
				Signal: req.Signal,
			}
			payload, err = json.Marshal(resp)
			if err != nil {
				slog.Error("fail to marshalling conversationSignal",
					"resp", resp)
				err = conn.Write(ctx, websocket.MessageText, []byte("fail to marshal"))
				if err != nil {
					slog.Error("fail to write payload",
						"err", err,
					)
				}
				return
			}
			err = to.Write(ctx, websocket.MessageText, payload)
			if err != nil {
				slog.Error("fail to write payload",
					"err", err,
				)
				return
			}
			continue
		}
		err = c.service.PublishConversationSignal(memberIdRaw, req.ToId, req.Signal)
		if err != nil {
			err = conn.Write(init, websocket.MessageText, []byte("fail to publish"))
			if err != nil {
				slog.Error("fail to write payload",
					"err", err,
				)
			}
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
