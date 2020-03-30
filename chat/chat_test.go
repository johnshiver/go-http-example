package chat

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/johnshiver/asapp_challenge/users"
	uTest "github.com/johnshiver/asapp_challenge/utils/testing"
)

func createTestUser(t *testing.T, tx *sqlx.Tx, username, password string) *users.User {
	userMan := users.UserManager{}
	user, err := userMan.Insert(tx, username, password)
	require.Nil(t, err)
	return user
}

func TestInsertMessageText(t *testing.T) {
	t.Parallel()
	var (
		username1      = "test-user"
		password1      = "my-password"
		username2      = "test-user2"
		password2      = "my-password"
		user1          *users.User
		user2          *users.User
		messageContent Content = map[string]interface{}{
			"type": "text",
			"text": "hey there",
		}
		chatMan Manager
	)

	ctx := context.Background()
	ctx, cancelTX := context.WithCancel(ctx)
	defer cancelTX()

	testDB := uTest.GetTestDB(t)
	tx, err := testDB.BeginTxx(ctx, nil)
	require.Nil(t, err)

	user1 = createTestUser(t, tx, username1, password1)
	user2 = createTestUser(t, tx, username2, password2)

	id, timestamp, err := chatMan.Insert(tx, user1.ID, user2.ID, messageContent)
	require.Nil(t, err)
	require.NotNil(t, id)
	require.NotNil(t, timestamp)
}

func TestInsertMessageImage(t *testing.T) {
	t.Parallel()
	var (
		username1      = "test-user"
		password1      = "my-password"
		username2      = "test-user2"
		password2      = "my-password"
		user1          *users.User
		user2          *users.User
		messageContent Content = map[string]interface{}{
			"type":   "image",
			"url":    "https://www.myimage.com/",
			"height": 1,
			"width":  1,
		}
		chatMan Manager
	)

	ctx := context.Background()
	ctx, cancelTX := context.WithCancel(ctx)
	defer cancelTX()

	testDB := uTest.GetTestDB(t)
	tx, err := testDB.BeginTxx(ctx, nil)
	require.Nil(t, err)

	user1 = createTestUser(t, tx, username1, password1)
	user2 = createTestUser(t, tx, username2, password2)

	id, timestamp, err := chatMan.Insert(tx, user1.ID, user2.ID, messageContent)
	require.Nil(t, err)
	require.NotNil(t, id)
	require.NotNil(t, timestamp)
}
func TestInsertMessageVideo(t *testing.T) {
	t.Parallel()
	var (
		username1      = "test-user"
		password1      = "my-password"
		username2      = "test-user2"
		password2      = "my-password"
		user1          *users.User
		user2          *users.User
		messageContent Content = map[string]interface{}{
			"type":   "video",
			"url":    "https://www.myvideo.com/",
			"source": "vimeo",
		}
		chatMan Manager
	)

	ctx := context.Background()
	ctx, cancelTX := context.WithCancel(ctx)
	defer cancelTX()

	testDB := uTest.GetTestDB(t)
	tx, err := testDB.BeginTxx(ctx, nil)
	require.Nil(t, err)

	user1 = createTestUser(t, tx, username1, password1)
	user2 = createTestUser(t, tx, username2, password2)

	id, timestamp, err := chatMan.Insert(tx, user1.ID, user2.ID, messageContent)
	require.Nil(t, err)
	require.NotNil(t, id)
	require.NotNil(t, timestamp)
}

func TestGetMessagesByRecipient(t *testing.T) {
	t.Parallel()
	var (
		username1   = "test-user"
		password1   = "my-password"
		username2   = "test-user1"
		password2   = "my-password"
		user1       *users.User
		user2       *users.User
		textMessage Content = map[string]interface{}{
			"type": "text",
			"text": "hey there",
		}
		imageMessage Content = map[string]interface{}{
			"type":   "image",
			"url":    "https://www.myimage.com/",
			"height": 1,
			"width":  1,
		}
		videoMessage Content = map[string]interface{}{
			"type":   "video",
			"url":    "https://www.myvideo.com/",
			"source": "vimeo",
		}
		chatMan Manager
		start   int64 = 0
		limit         = 100
	)

	ctx := context.Background()
	ctx, cancelTX := context.WithCancel(ctx)
	defer cancelTX()

	testDB := uTest.GetTestDB(t)
	tx, err := testDB.BeginTxx(ctx, nil)
	require.Nil(t, err)

	user1 = createTestUser(t, tx, username1, password1)
	user2 = createTestUser(t, tx, username2, password2)

	var content Content
	for i := 0; i < limit*2; i++ {
		if i%2 == 0 {
			content = textMessage
		} else if i%3 == 0 {
			content = imageMessage
		} else {
			content = videoMessage
		}
		id, timestamp, err := chatMan.Insert(tx, user1.ID, user2.ID, content)
		require.Nil(t, err)
		require.NotNil(t, id)
		require.NotNil(t, timestamp)
	}

	messageBatch1, err := chatMan.GetMessagesByRecipient(tx, user1.ID, user2.ID, start, limit)
	require.Nil(t, err)
	require.Len(t, messageBatch1, limit)

	start = messageBatch1[len(messageBatch1)-1].ID
	messageBatch2, err := chatMan.GetMessagesByRecipient(tx, user1.ID, user2.ID, start+1, limit)
	require.Nil(t, err)
	require.Len(t, messageBatch2, limit)
	require.Greater(t, messageBatch2[0].ID, messageBatch1[limit-1].ID)

}
