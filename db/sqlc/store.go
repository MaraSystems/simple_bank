package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("tx err: %w, rb err: %w", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:from_account_id`
	ToAccountID   int64 `json:to_account_id`
	Amount        int64 `json:amount`
}

type TransferTxResult struct {
	Transfer    Transfer `json:transfer`
	FromAccount Account  `json:from_account`
	ToAccount   Account  `json:to_account`
	FromEntry   Entry    `json:from_entry`
	ToEntry     Entry    `json:to_entry`
}

var txKey = struct{}{}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(qtx *Queries) error {
		var err error

		result.Transfer, err = qtx.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: sql.NullInt64{Int64: arg.FromAccountID, Valid: true},
			ToAccountID:   sql.NullInt64{Int64: arg.ToAccountID, Valid: true},
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		result.FromEntry, err = qtx.CreateEntry(ctx, CreateEntryParams{
			AccountID: sql.NullInt64{Int64: arg.FromAccountID, Valid: true},
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = qtx.CreateEntry(ctx, CreateEntryParams{
			AccountID: sql.NullInt64{Int64: arg.ToAccountID, Valid: true},
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		// if arg.FromAccountID > arg.ToAccountID {
		// 	result.FromAccount, err = qtx.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		// 		ID:     arg.FromAccountID,
		// 		Amount: -arg.Amount,
		// 	})
		// } else {
		// 	result.ToAccount, err = qtx.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		// 		ID:     arg.ToAccountID,
		// 		Amount: arg.Amount,
		// 	})
		// }

		if arg.FromAccountID > arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, qtx, arg.FromAccountID, arg.ToAccountID, -arg.Amount, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, qtx, arg.ToAccountID, arg.FromAccountID, arg.Amount, -arg.Amount)
		}

		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	fromAccountID int64,
	toAccountID int64,
	fromAmount int64,
	toAmount int64,
) (fromAccount, toAccount Account, err error) {
	fromAccount, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID:     fromAccountID,
		Amount: fromAmount,
	})
	if err != nil {
		return
	}

	toAccount, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID:     toAccountID,
		Amount: toAmount,
	})
	return
}
