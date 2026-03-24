package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)

	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	err = CheckPassword(hashedPassword, password)
	require.NoError(t, err)

	invalidPassword := fmt.Sprintf("%va", password)
	err = CheckPassword(hashedPassword, invalidPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	anotherHahsedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, anotherHahsedPassword)
	require.NotEqual(t, hashedPassword, anotherHahsedPassword)
}
