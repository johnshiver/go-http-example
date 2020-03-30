package users

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/require"

	uTest "github.com/johnshiver/asapp_challenge/utils/testing"
)

func TestLogin(t *testing.T) {
	t.Parallel()

	var (
		username = "string"
		password = "pa$$word"
		reqBody  = userRequest{
			Username: username,
			Password: password,
		}
		testUserID int64 = 1
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mUserServ := NewMockUserService(ctrl)
	mUserServ.EXPECT().
		Authenticate(gomock.Any(), username, password).
		Return(testUserID, nil)

	testDB := uTest.GetTestDB(t)
	env := env{
		db:          testDB,
		userService: mUserServ,
	}

	router := httprouter.New()
	env.addRoutes(router)

	reqJSON, err := json.Marshal(&reqBody)
	require.Nil(t, err)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(reqJSON))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, w.Code, http.StatusOK)

	var resp loginResponse
	err = json.NewDecoder(w.Body).Decode(&resp)
	require.Nil(t, err)
	require.NotNil(t, resp.Token)
	require.Equal(t, resp.ID, testUserID)
}
