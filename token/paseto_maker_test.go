package token

import (
	"github.com/WanCodeBase/GinModule/util"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	assert.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := time.Now().Add(time.Minute)

	token, err := maker.CreateToken(username, duration)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	assert.NoError(t, err)
	assert.NotEmpty(t, payload)
	assert.NotZero(t, payload.ID)
	assert.Equal(t, username, payload.Username)
	assert.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	assert.WithinDuration(t, expiredAt, payload.ExpireAt, time.Second)
}

func TestExpiredPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	assert.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	token, err := maker.CreateToken(username, -duration)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	_, err = maker.VerifyToken(token)
	assert.ErrorIs(t, ErrExpireToken, err)
}
