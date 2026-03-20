package db

import (
	"database/sql"
	"testing"
	"time"

	"github.com/odekodu/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) Transfer {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	arg := CreateTransferParams{
		FromAccountID: sql.NullInt64{Int64: fromAccount.ID, Valid: true},
		ToAccountID:   sql.NullInt64{Int64: toAccount.ID, Valid: true},
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(t.Context(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	randomTransfer := createRandomTransfer(t)

	transfer, err := testQueries.GetTransfer(t.Context(), randomTransfer.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, randomTransfer.FromAccountID, transfer.FromAccountID)
	require.Equal(t, randomTransfer.ToAccountID, transfer.ToAccountID)
	require.Equal(t, randomTransfer.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
}

func TestListTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomTransfer(t)
	}

	arg := ListTransfersParams{
		Limit:  5,
		Offset: 0,
	}

	entries, err := testQueries.ListTransfers(t.Context(), arg)

	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, transfer := range entries {
		require.NotEmpty(t, transfer)
	}
}

func TestUpdateTransfer(t *testing.T) {
	randomTransfer := createRandomTransfer(t)
	arg := UpdateTransferParams{
		ID:     randomTransfer.ID,
		Amount: util.RandomMoney(),
	}

	transfer, err := testQueries.UpdateTransfer(t.Context(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, randomTransfer.FromAccountID, transfer.FromAccountID)
	require.Equal(t, randomTransfer.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, randomTransfer.ID)
	require.NotZero(t, randomTransfer.CreatedAt)
	require.WithinDuration(t, randomTransfer.CreatedAt.Time, transfer.CreatedAt.Time, time.Second)
}

func TestDeleteTransfer(t *testing.T) {
	randomTransfer := createRandomTransfer(t)
	err := testQueries.DeleteTransfer(t.Context(), randomTransfer.ID)
	require.NoError(t, err)

	transfer, err := testQueries.GetTransfer(t.Context(), randomTransfer.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, transfer)
}
