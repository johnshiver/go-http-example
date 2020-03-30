package users

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/johnshiver/asapp_challenge/utils"
)

type userRequest struct {
	Username string `json:"username" validate:"required,gte=4,lte=25,alphanumunicode"`
	Password string `json:"password" validate:"required,gte=8,lte=100"`
}

func (ur *userRequest) Validate() error {
	v := utils.GetValidator()
	err := v.Struct(ur)
	if err != nil {
		return err
	}
	return nil
}

func (e *env) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancelTX := context.WithCancel(r.Context())
	defer cancelTX()
	tx, err := e.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Printf("starting transaction: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	var userPayload userRequest
	if err = json.NewDecoder(r.Body).Decode(&userPayload); err != nil {
		log.Printf("decoding user payload: %v", err)
		http.Error(w, "bad payload", http.StatusBadRequest)
		return
	}

	if errors := userPayload.Validate(); errors != nil {
		log.Printf("validating user payload: %v", errors)
		http.Error(w, errors.Error(), http.StatusBadRequest)
		return
	}

	user, err := e.userService.Insert(tx, userPayload.Username, userPayload.Password)
	if err != nil {
		log.Printf("inserting user row: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("commiting user row: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	rJSON, err := json.Marshal(map[string]int64{"id": user.ID})
	if err != nil {
		log.Printf("creating response json: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Context-Type", "application/json")
	_, _ = w.Write(rJSON)
}
