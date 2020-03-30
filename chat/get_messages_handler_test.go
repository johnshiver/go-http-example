package chat

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/require"

	"github.com/johnshiver/asapp_challenge/utils"
	uTest "github.com/johnshiver/asapp_challenge/utils/testing"
)

func TestGetMessages(t *testing.T) {
	var (
		senderId    int64 = 1
		recipientId int64 = 2
		start       int64 = 0
		limit             = 100
		queryParams       = map[string]string{
			"recipient": "2",
			"start":     "0",
		}
		messages = []*Message{
			{
				ID:          0,
				RecipientID: 0,
				SenderID:    0,
				Content:     nil,
			},
			{
				ID:          0,
				RecipientID: 0,
				SenderID:    0,
				Content:     nil,
			},
		}
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mChatServ := NewMockService(ctrl)
	mChatServ.EXPECT().
		GetMessagesByRecipient(gomock.Any(), senderId, recipientId, start, limit).
		Return(messages, nil)

	testDB := uTest.GetTestDB(t)
	env := env{
		db:          testDB,
		chatService: mChatServ,
	}
	router := httprouter.New()
	env.addRoutes(router)

	req := httptest.NewRequest("GET", "/messages", nil)
	q := req.URL.Query()
	for k, v := range queryParams {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	// add token to request header
	token, err := utils.CreateToken(senderId)
	require.Nil(t, err)
	req.Header.Add("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, w.Code, http.StatusOK)

	var resp map[string]interface{}
	err = json.NewDecoder(w.Body).Decode(&resp)
	require.Nil(t, err)
	require.NotNil(t, resp)
}
