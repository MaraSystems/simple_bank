package db

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	sendAccount := createRandomAccount(t)
	receiveAccount := createRandomAccount(t)
	store := NewStore(testDb)

	amount := int64(10)
	n := 5

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(t.Context(), TransferTxParams{
				FromAccountID: sendAccount.ID,
				ToAccountID:   receiveAccount.ID,
				Amount:        int64(amount),
			})

			errs <- err
			results <- result
		}()
	}

	fmt.Println(">> Before: ", sendAccount.Balance, receiveAccount.Balance)
	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		result := <-results

		require.NoError(t, err)
		require.NotEmpty(t, result)

		// Check transfer
		tranfer := result.Transfer
		require.NotEmpty(t, tranfer)
		require.Equal(t, sendAccount.ID, tranfer.FromAccountID.Int64)
		require.Equal(t, receiveAccount.ID, tranfer.ToAccountID.Int64)
		require.Equal(t, amount, tranfer.Amount)
		require.NotZero(t, tranfer.CreatedAt)
		require.NotZero(t, tranfer.ID)

		_, err = store.GetTransfer(t.Context(), tranfer.ID)
		require.NoError(t, err)

		// Check for entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromEntry.AccountID.Int64, sendAccount.ID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.CreatedAt)
		require.NotZero(t, fromEntry.ID)

		_, err = store.GetEntry(t.Context(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, toEntry.AccountID.Int64, receiveAccount.ID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.CreatedAt)
		require.NotZero(t, toEntry.ID)

		_, err = store.GetEntry(t.Context(), toEntry.ID)
		require.NoError(t, err)

		// Check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, sendAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, receiveAccount.ID)

		fromDiff := sendAccount.Balance - fromAccount.Balance
		toDiff := toAccount.Balance - receiveAccount.Balance
		require.Equal(t, fromDiff, toDiff)
		require.True(t, fromDiff > 0)
		require.True(t, fromDiff%amount == 0) // (1, 2, 3, ...) * amount

		fmt.Println(">> Tx: ", fromAccount.Balance, toAccount.Balance)

		k := int(fromDiff / amount)
		fmt.Println(">> k", k)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updatedSendAccount, err := store.GetAccount(t.Context(), sendAccount.ID)
	require.NoError(t, err)
	require.Equal(t, updatedSendAccount.Balance, sendAccount.Balance-(int64(n)*amount))

	updatedReceiveAccount, err := store.GetAccount(t.Context(), receiveAccount.ID)
	require.NoError(t, err)
	require.Equal(t, updatedReceiveAccount.Balance, receiveAccount.Balance+(int64(n)*amount))

	fmt.Println(">> After: ", updatedSendAccount.Balance, updatedReceiveAccount.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	sendAccount := createRandomAccount(t)
	receiveAccount := createRandomAccount(t)
	store := NewStore(testDb)

	amount := int64(10)
	n := 10

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID, toAccountID := sendAccount.ID, receiveAccount.ID
		if i%2 == 1 {
			fromAccountID, toAccountID = toAccountID, fromAccountID
		}

		go func() {
			_, err := store.TransferTx(t.Context(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        int64(amount),
			})

			errs <- err
		}()
	}

	fmt.Println(">> Before: ", sendAccount.Balance, receiveAccount.Balance)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	updatedSendAccount, err := store.GetAccount(t.Context(), sendAccount.ID)
	require.NoError(t, err)
	require.Equal(t, updatedSendAccount.Balance, sendAccount.Balance)

	updatedReceiveAccount, err := store.GetAccount(t.Context(), receiveAccount.ID)
	require.NoError(t, err)
	require.Equal(t, updatedReceiveAccount.Balance, receiveAccount.Balance)

	fmt.Println(">> After: ", updatedSendAccount.Balance, updatedReceiveAccount.Balance)
}
