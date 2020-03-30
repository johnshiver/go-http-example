package chat

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/johnshiver/asapp_challenge/middleware"
)

type getMessagesResp struct {
	Messages []*Message `json:"messages"`
}

func (e *env) GetMessages(w http.ResponseWriter, r *http.Request) {
	ctx, cancelTX := context.WithCancel(r.Context())
	defer cancelTX()
	tx, err := e.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Printf("starting transaction: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	requestUserID, err := middleware.GetUserIDFromRequest(r)
	if err != nil {
		log.Printf("no user id in request context")
		http.Error(w, "unauthorized request", http.StatusForbidden)
		return
	}

	// get query params
	query := r.URL.Query()
	recipientRaw, ok := query["recipient"]
	if !ok || len(recipientRaw) < 1 {
		log.Printf("no recipient query param")
		http.Error(w, "recipient query param required", http.StatusBadRequest)
		return
	}
	startRaw, ok := query["start"]
	if !ok || len(startRaw) < 1 {
		log.Printf("no start query param")
		http.Error(w, "start query param required", http.StatusBadRequest)
		return
	}

	var limitRaw []string
	limitRaw, ok = query["limit"]
	if !ok || len(limitRaw) < 1 {
		limitRaw = []string{"100"}
	}

	recipientID, err := strconv.ParseInt(recipientRaw[0], 10, 64)
	if err != nil {
		log.Printf("bad recipient id %v", err)
		http.Error(w, "recipient query param must be int", http.StatusBadRequest)
		return
	}
	start, err := strconv.ParseInt(startRaw[0], 10, 64)
	if err != nil {
		log.Printf("bad start %v", err)
		http.Error(w, "start query param must be int", http.StatusBadRequest)
		return
	}
	limit, err := strconv.Atoi(limitRaw[0])
	if err != nil {
		log.Printf("bad limit %v", err)
		http.Error(w, "limit query param must be int", http.StatusBadRequest)
		return
	}

	messages, err := e.chatService.GetMessagesByRecipient(tx, requestUserID, recipientID, start, limit)
	if err != nil {
		log.Printf("getting messages from db: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	rJSON, err := json.Marshal(&getMessagesResp{messages})
	if err != nil {
		log.Printf("creating response json: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Context-Type", "application/json")
	_, _ = w.Write(rJSON)
}
