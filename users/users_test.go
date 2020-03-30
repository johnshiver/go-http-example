package users

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	uTest "github.com/johnshiver/asapp_challenge/utils/testing"
)

func TestInsertUser(t *testing.T) {
	t.Parallel()
	var (
		username = "test-user"
		password = "my-password"
	)

	ctx := context.Background()
	ctx, cancelTX := context.WithCancel(ctx)
	defer cancelTX()

	testDB := uTest.GetTestDB(t)
	tx, err := testDB.BeginTxx(ctx, nil)
	require.Nil(t, err)

	userMang := UserManager{}
	testUser, err := userMang.Insert(tx, username, password)
	require.Nil(t, err)
	require.Equal(t, username, testUser.Username)
}

func TestInsertDuplicateUser(t *testing.T) {
	t.Parallel()
	var (
		username  = "test-user"
		password  = "my-password"
		username2 = "test-user"
		password2 = "my-password"
	)

	ctx := context.Background()
	ctx, cancelTX := context.WithCancel(ctx)
	defer cancelTX()

	testDB := uTest.GetTestDB(t)
	tx, err := testDB.BeginTxx(ctx, nil)
	require.Nil(t, err)

	userMang := UserManager{}
	testUser, err := userMang.Insert(tx, username, password)
	require.Nil(t, err)
	require.Equal(t, username, testUser.Username)

	testUser2, err := userMang.Insert(tx, username2, password2)
	require.Nil(t, testUser2)
	require.NotNil(t, err)
}

func TestGetByID(t *testing.T) {
	t.Parallel()
	var (
		username = "test-user"
		password = "my-password"
	)

	ctx := context.Background()
	ctx, cancelTX := context.WithCancel(ctx)
	defer cancelTX()

	testDB := uTest.GetTestDB(t)
	tx, err := testDB.BeginTxx(ctx, nil)
	require.Nil(t, err)

	userMang := UserManager{}
	testUser, err := userMang.Insert(tx, username, password)
	require.Nil(t, err)
	require.Equal(t, username, testUser.Username)

	shouldBeSameUser, err := userMang.GetByID(tx, testUser.ID)
	require.Nil(t, err)
	require.Equal(t, testUser.ID, shouldBeSameUser.ID)
}

func TestAuthenticate(t *testing.T) {
	t.Parallel()
	var (
		username = "test-user"
		password = "my-password"
	)

	ctx := context.Background()
	ctx, cancelTX := context.WithCancel(ctx)
	defer cancelTX()

	testDB := uTest.GetTestDB(t)
	tx, err := testDB.BeginTxx(ctx, nil)
	require.Nil(t, err)

	userMang := UserManager{}
	testUser, err := userMang.Insert(tx, username, password)
	require.Nil(t, err)
	require.Equal(t, username, testUser.Username)

	userID, err := userMang.Authenticate(tx, username, password)
	require.Nil(t, err)
	require.Equal(t, testUser.ID, userID)
}
