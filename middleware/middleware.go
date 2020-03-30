package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/justinas/alice"

	"github.com/johnshiver/asapp_challenge/utils"
)

type contextKey string

var (
	ContextKeyUserID = contextKey("userID")
	Standard         = alice.New(
		RecoverPanic,
		LogRequest,
	)
)

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s - %s %s %s",
			r.RemoteAddr,
			r.Proto,
			r.Method,
			r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func RecoverPanic(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				log.Printf("recovered: %v", err)
				w.Header().Set("Connection", "close")
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func AuthenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Header should come in the form: Authorization: Bearer $token
		reqToken := r.Header.Get("Authorization")
		if len(reqToken) < 2 {
			http.Error(w, "no Authorization header", http.StatusForbidden)
			return
		}
		splitToken := strings.Split(reqToken, "Bearer ")
		if len(splitToken) != 2 {
			http.Error(w, "no token", http.StatusForbidden)
			return
		}
		rawToken := splitToken[1]
		if len(rawToken) < 2 {
			http.Error(w, "bad token", http.StatusForbidden)
			return
		}

		claims, err := utils.VerifyToken(rawToken)
		if err != nil {
			log.Printf("authenticating user: %v", err)
			http.Error(w, "bad token", http.StatusForbidden)
			return
		}
		ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
