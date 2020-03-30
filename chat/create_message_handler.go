package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/johnshiver/asapp_challenge/middleware"
	"github.com/johnshiver/asapp_challenge/utils"
)

type messagePayload struct {
	SenderID    int64   `json:"sender" validate:"required,numeric,gte=1,nefield=RecipientID"`
	RecipientID int64   `json:"recipient" validate:"required,numeric,gte=1,nefield=SenderID"`
	Content     Content `json:"content" validate:"required"`
}

func (mp *messagePayload) Validate() error {
	v := utils.GetValidator()
	err := v.Struct(mp)
	if err != nil {
		return fmt.Errorf("validating message payload: %v", err)
	}
	err = mp.Content.Validate()
	if err != nil {
		return err
	}
	return nil
}

type createMsgResp struct {
	ID        int64          `json:"id"`
	Timestamp utils.JSONTime `json:"timestamp"`
}

func (e *env) CreateMessage(w http.ResponseWriter, r *http.Request) {
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

	var payload messagePayload
	if err = json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("decoding message payload: %v", err)
		http.Error(w, "bad payload", http.StatusBadRequest)
		return
	}

	if errors := payload.Validate(); errors != nil {
		log.Printf("validating message payload: %v", errors)
		http.Error(w, errors.Error(), http.StatusBadRequest)
		return
	}

	if payload.SenderID != requestUserID {
		log.Printf("request user %d tried to send message on behalf of %d", requestUserID, payload.SenderID)
		http.Error(w, "unauthorized", http.StatusForbidden)
		return
	}

	messageID, messageTimestamp, err := e.chatService.Insert(tx, payload.SenderID, payload.RecipientID, payload.Content)
	if err != nil {
		log.Printf("inserting message row: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("commiting message row: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := createMsgResp{
		ID:        messageID,
		Timestamp: utils.JSONTime(messageTimestamp),
	}
	rJSON, err := json.Marshal(&resp)
	if err != nil {
		log.Printf("creating response json: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Context-Type", "application/json")
	_, _ = w.Write(rJSON)
}
