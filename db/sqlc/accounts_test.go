package db

import (
	"context"
	"github.com/WanCodeBase/GinModule/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func _createAccount(t *testing.T) Account {
	user := _createUser(t)
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	assert.NoError(t, err)
	assert.NotEmpty(t, account)

	assert.Equal(t, arg.Owner, account.Owner)
	assert.Equal(t, arg.Balance, account.Balance)
	assert.Equal(t, arg.Currency, account.Currency)

	assert.NotZero(t, account.ID)
	assert.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	_createAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := _createAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	assert.NoError(t, err)
	assert.NotEmpty(t, account2)
	assert.Equal(t, account1, account2)
}

func TestUpdateAccount(t *testing.T) {
	account1 := _createAccount(t)
	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}
	account2, err := testQueries.UpdateAccount(context.Background(), arg)

	assert.NoError(t, err)
	assert.NotEmpty(t, account2)

	assert.Equal(t, arg.Balance, account2.Balance)
	assert.Equal(t, account1.ID, account2.ID)
	assert.Equal(t, account1.Owner, account2.Owner)
}

func TestDeleteAccount(t *testing.T) {
	account1 := _createAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)

	assert.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	assert.Error(t, err)
	assert.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	var lastUsername string
	for i := 0; i < 6; i++ {
		account := _createAccount(t)
		lastUsername = account.Owner
	}

	arg := ListAccountsParams{
		Owner:  lastUsername,
		Limit:  5,
		Offset: 0,
	}
	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(accounts))
}
