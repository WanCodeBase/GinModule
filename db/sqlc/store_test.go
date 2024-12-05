package db

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStore_TransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := _createAccount(t)
	account2 := _createAccount(t)

	n := 5
	amount := int64(10)

	// channel
	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		assert.NoError(t, err)

		result := <-results
		assert.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		assert.NotEmpty(t, transfer)
		assert.Equal(t, account1.ID, transfer.FromAccountID)
		assert.Equal(t, account2.ID, transfer.ToAccountID)
		assert.Equal(t, amount, transfer.Amount)
		assert.NotZero(t, transfer.ID)
		assert.NotZero(t, transfer.CreatedAt)

		_, err = testQueries.GetTransfer(context.Background(), transfer.ID)
		assert.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		assert.NotEmpty(t, fromEntry)
		assert.NotZero(t, fromEntry.ID)
		assert.Equal(t, account1.ID, fromEntry.AccountID)
		assert.Equal(t, -amount, fromEntry.Amount)
		assert.NotZero(t, fromEntry.CreatedAt)

		toEntry := result.ToEntry
		assert.NotEmpty(t, toEntry)
		assert.NotZero(t, toEntry.ID)
		assert.Equal(t, account2.ID, toEntry.AccountID)
		assert.Equal(t, amount, toEntry.Amount)
		assert.NotZero(t, toEntry.CreatedAt)

		// check accounts
		fromAccount := result.FromAccount
		assert.NotEmpty(t, fromAccount)
		assert.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		assert.NotEmpty(t, toAccount)
		assert.Equal(t, account2.ID, toAccount.ID)

		// check balance
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		assert.Equal(t, diff1, diff2)
		assert.True(t, diff1 > 0)
		assert.True(t, diff1%amount == 0)
	}

	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	assert.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	assert.NoError(t, err)

	assert.Equal(t, account1.Balance-int64(n)*amount, updateAccount1.Balance)
	assert.Equal(t, account2.Balance+int64(n)*amount, updateAccount2.Balance)
}

func TestStore_TransferTxDeadlockTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := _createAccount(t)
	account2 := _createAccount(t)

	n := 10
	amount := int64(10)

	// channel
	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccount, toAccount := account1, account2

		if i%2 == 0 {
			fromAccount, toAccount = account2, account1
		}
		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        amount,
			})
			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		assert.NoError(t, err)

	}

	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	assert.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	assert.NoError(t, err)

	assert.Equal(t, account1.Balance, updateAccount1.Balance)
	assert.Equal(t, account2.Balance, updateAccount2.Balance)
}
