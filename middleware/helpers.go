package middleware

import (
	"fmt"
	"net/http"
)

func GetUserIDFromRequest(request *http.Request) (int64, error) {
	userID, ok := request.Context().Value(ContextKeyUserID).(int64)
	if !ok {
		err := fmt.Errorf("error setting userid in context")
		return 0, err
	}
	return userID, nil
}
