package users

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/johnshiver/asapp_challenge/utils"
)

type loginResponse struct {
	ID    int64  `json:"id"`
	Token string `json:"token"`
}

func (e *env) Login(w http.ResponseWriter, r *http.Request) {
	ctx, cancelTX := context.WithCancel(r.Context())
	defer cancelTX()
	tx, err := e.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Printf("starting transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var userPayload userRequest
	if err = json.NewDecoder(r.Body).Decode(&userPayload); err != nil {
		log.Printf("decoding user payload: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// dont need to validate in this case
	userID, err := e.userService.Authenticate(tx, userPayload.Username, userPayload.Password)
	if err != nil {
		if err == ErrNoRecord {
			log.Printf("no record found for %s", userPayload.Username)
			http.Error(w, "bad credentials", http.StatusForbidden)
			return
		}
		if err == ErrInvalidCredentials {
			log.Printf("bad creds for %s", userPayload.Username)
			http.Error(w, "bad credentials", http.StatusForbidden)
			return
		}
		log.Printf("error while authenticating user: %s -> %v", userPayload.Username, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	token, err := utils.CreateToken(userID)
	if err != nil {
		log.Printf("creating token %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	loginResp := loginResponse{
		ID:    userID,
		Token: token,
	}

	rJSON, err := json.Marshal(&loginResp)
	if err != nil {
		log.Printf("creating response json: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Context-Type", "application/json")
	_, _ = w.Write(rJSON)
}
