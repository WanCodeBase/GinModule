package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password1 := RandomString(6)
	hashedPassword, err := HashedPassword(password1)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)

	err = checkPassword(password1, hashedPassword)
	assert.NoError(t, err)

	password2 := RandomString(6)
	err = checkPassword(password2, hashedPassword)
	assert.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
