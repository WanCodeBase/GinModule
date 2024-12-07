package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// SQLStore provides all functions to execute db queries & transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(queries *Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil) // default level: RC
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		// rollback
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %s, rb error: %s", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

/*
TransferTx performs a money transfer one account to another
1. creates a new transfers
2. add account entries
3. and update accounts' balance
within a single database transaction
*/
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(queries *Queries) error {
		transfer, err := queries.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		fromEntry, err := queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    arg.Amount * (-1),
		})
		if err != nil {
			return err
		}

		toEntry, err := queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// update accounts' balance
		// 预防死锁：确保获取锁的顺序是一致的 （eg.总是id小的对象先获取锁）
		// Prevent deadlocks: Ensure that locks are acquired in the same order
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = store.addMoney(ctx, queries,
				AddAccountBalanceParams{
					ID:     arg.FromAccountID,
					Amount: -1 * arg.Amount,
				}, AddAccountBalanceParams{
					ID:     arg.ToAccountID,
					Amount: arg.Amount,
				})
		} else {
			result.ToAccount, result.FromAccount, err = store.addMoney(ctx, queries,
				AddAccountBalanceParams{
					ID:     arg.ToAccountID,
					Amount: arg.Amount,
				}, AddAccountBalanceParams{
					ID:     arg.FromAccountID,
					Amount: -1 * arg.Amount,
				})
		}

		result.Transfer = transfer
		result.FromEntry = fromEntry
		result.ToEntry = toEntry
		return nil
	})
	if err != nil {
		log.Fatal("transfer err:", err)
	}

	return result, err
}

func (store *SQLStore) addMoney(ctx context.Context, q *Queries, param1, param2 AddAccountBalanceParams) (account1, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, param1)
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, param2)
	return
}
