package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/WanCodeBase/GinModule/util"
	"github.com/stretchr/testify/assert"
)

func _createUser(t *testing.T) User {
	password, err := util.HashedPassword(util.RandomString(8))
	assert.NoError(t, err)
	
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: password,
		FullName:       util.RandomOwner(),
		Email:          fmt.Sprintf("%s@example.com", util.RandomOwner()),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	assert.NoError(t, err)
	assert.NotEmpty(t, user)

	assert.Equal(t, arg.Username, user.Username)
	assert.Equal(t, arg.HashedPassword, user.HashedPassword)
	assert.Equal(t, arg.FullName, user.FullName)
	assert.Equal(t, arg.Email, user.Email)

	assert.NotZero(t, user.CreatedAt)
	assert.Zero(t, user.PasswordChangedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	_createUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := _createUser(t)

	user2, err := testQueries.GetUser(context.Background(), user1.Username)

	assert.NoError(t, err)
	assert.NotEmpty(t, user2)
	assert.Equal(t, user1, user2)
}
