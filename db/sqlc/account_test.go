package db

import (
	"database/sql"
	"testing"
	"time"

	"github.com/MaraSystems/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(t.Context(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	randomAccount := createRandomAccount(t)
	account, err := testQueries.GetAccount(t.Context(), randomAccount.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, randomAccount.Owner, account.Owner)
	require.Equal(t, randomAccount.Balance, account.Balance)
	require.Equal(t, randomAccount.Currency, account.Currency)

	require.NotZero(t, randomAccount.ID)
	require.NotZero(t, randomAccount.CreatedAt)
	require.WithinDuration(t, randomAccount.CreatedAt.Time, account.CreatedAt.Time, time.Second)
}

func TestUpdateAccountBalance(t *testing.T) {
	randomAccount := createRandomAccount(t)
	arg := UpdateAccountBalanceParams{
		ID:     randomAccount.ID,
		Amount: util.RandomMoney(),
	}

	account, err := testQueries.UpdateAccountBalance(t.Context(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, randomAccount.Owner, account.Owner)
	require.Equal(t, randomAccount.Balance+arg.Amount, account.Balance)
	require.Equal(t, randomAccount.Currency, account.Currency)

	require.NotZero(t, randomAccount.ID)
	require.NotZero(t, randomAccount.CreatedAt)
	require.WithinDuration(t, randomAccount.CreatedAt.Time, account.CreatedAt.Time, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	randomAccount := createRandomAccount(t)
	err := testQueries.DeleteAccount(t.Context(), randomAccount.ID)
	require.NoError(t, err)

	account, err := testQueries.GetAccount(t.Context(), randomAccount.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(t.Context(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, acc := range accounts {
		require.NotEmpty(t, acc)
		require.Equal(t, lastAccount.Owner, acc.Owner)
	}
}
