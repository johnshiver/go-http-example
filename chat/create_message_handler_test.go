package chat

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/require"

	"github.com/johnshiver/asapp_challenge/utils"
	uTest "github.com/johnshiver/asapp_challenge/utils/testing"
)

func TestCreateMessageText(t *testing.T) {
	var (
		senderId    int64 = 1
		recipientId int64 = 2
		reqBody           = messagePayload{
			SenderID:    senderId,
			RecipientID: recipientId,
			Content: map[string]interface{}{
				"type": "text",
				"text": "hey there",
			},
		}
		messageID        int64 = 1
		messageTimeStamp       = time.Now()
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mChatServ := NewMockService(ctrl)
	mChatServ.EXPECT().
		Insert(gomock.Any(), reqBody.SenderID, reqBody.RecipientID, reqBody.Content).
		Return(messageID, messageTimeStamp, nil)

	testDB := uTest.GetTestDB(t)
	env := env{
		db:          testDB,
		chatService: mChatServ,
	}

	router := httprouter.New()
	env.addRoutes(router)

	reqJSON, err := json.Marshal(&reqBody)
	require.Nil(t, err)
	req := httptest.NewRequest("POST", "/messages", bytes.NewBuffer(reqJSON))

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

func TestCreateMessageImage(t *testing.T) {
	var (
		senderId    int64 = 1
		recipientId int64 = 2
		reqBody           = messagePayload{
			SenderID:    senderId,
			RecipientID: recipientId,
			Content: map[string]interface{}{
				"type":   "image",
				"url":    "https://www.imageurl.com",
				"height": 1,
				"width":  2,
			},
		}
		messageID        int64 = 1
		messageTimeStamp       = time.Now()
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mChatServ := NewMockService(ctrl)
	// for now, gomock.Any() for content
	mChatServ.EXPECT().
		Insert(gomock.Any(), reqBody.SenderID, reqBody.RecipientID, gomock.Any()).
		Return(messageID, messageTimeStamp, nil)

	testDB := uTest.GetTestDB(t)
	env := env{
		db:          testDB,
		chatService: mChatServ,
	}

	router := httprouter.New()
	env.addRoutes(router)

	reqJSON, err := json.Marshal(&reqBody)
	require.Nil(t, err)
	req := httptest.NewRequest("POST", "/messages", bytes.NewBuffer(reqJSON))

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

func TestCreateMessageVideo(t *testing.T) {
	var (
		senderId    int64 = 1
		recipientId int64 = 2
		reqBody           = messagePayload{
			SenderID:    senderId,
			RecipientID: recipientId,
			Content: map[string]interface{}{
				"type":   "video",
				"url":    "https://www.imageurl.com",
				"source": "youtube",
			},
		}
		messageID        int64 = 1
		messageTimeStamp       = time.Now()
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mChatServ := NewMockService(ctrl)
	// for now, gomock.Any() for content
	mChatServ.EXPECT().
		Insert(gomock.Any(), reqBody.SenderID, reqBody.RecipientID, gomock.Any()).
		Return(messageID, messageTimeStamp, nil)

	testDB := uTest.GetTestDB(t)
	env := env{
		db:          testDB,
		chatService: mChatServ,
	}

	router := httprouter.New()
	env.addRoutes(router)

	reqJSON, err := json.Marshal(&reqBody)
	require.Nil(t, err)
	req := httptest.NewRequest("POST", "/messages", bytes.NewBuffer(reqJSON))

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
