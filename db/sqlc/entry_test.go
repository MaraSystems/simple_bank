package db

import (
	"database/sql"
	"testing"
	"time"

	"github.com/odekodu/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {
	account := createRandomAccount(t)
	arg := CreateEntryParams{
		AccountID: sql.NullInt64{Int64: account.ID, Valid: true},
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(t.Context(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	randomEntry := createRandomEntry(t)

	entry, err := testQueries.GetEntry(t.Context(), randomEntry.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, randomEntry.AccountID, entry.AccountID)
	require.Equal(t, randomEntry.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
}

func TestListEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomEntry(t)
	}

	arg := ListEntriesParams{
		Limit:  5,
		Offset: 0,
	}

	entries, err := testQueries.ListEntries(t.Context(), arg)

	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}

func TestUpdateEntry(t *testing.T) {
	randomEntry := createRandomEntry(t)
	arg := UpdateEntryParams{
		ID:     randomEntry.ID,
		Amount: util.RandomMoney(),
	}

	entry, err := testQueries.UpdateEntry(t.Context(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, randomEntry.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, randomEntry.ID)
	require.NotZero(t, randomEntry.CreatedAt)
	require.WithinDuration(t, randomEntry.CreatedAt.Time, entry.CreatedAt.Time, time.Second)
}

func TestDeleteEntry(t *testing.T) {
	randomEntry := createRandomEntry(t)
	err := testQueries.DeleteEntry(t.Context(), randomEntry.ID)
	require.NoError(t, err)

	entry, err := testQueries.GetEntry(t.Context(), randomEntry.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entry)
}
