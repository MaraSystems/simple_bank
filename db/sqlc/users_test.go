package db

import (
	"fmt"
	"testing"

	"github.com/MaraSystems/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:       util.RandomString(5),
		HashedPassword: "secret",
		FullName:       fmt.Sprintf("%v %v", util.RandomString(5), util.RandomString(5)),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(t.Context(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordUpdatedAt.IsZero())
	require.NotEmpty(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	randomUser := createRandomUser(t)

	user, err := testQueries.GetUser(t.Context(), randomUser.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, randomUser.Username, user.Username)
	require.Equal(t, randomUser.HashedPassword, user.HashedPassword)
	require.Equal(t, randomUser.FullName, user.FullName)
	require.Equal(t, randomUser.Email, user.Email)

	require.True(t, user.PasswordUpdatedAt.IsZero())
	require.NotEmpty(t, user.CreatedAt)
}
