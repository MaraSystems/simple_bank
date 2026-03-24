package api

import (
	"os"
	"testing"

	db "github.com/MaraSystems/simple_bank/db/sqlc"
	"github.com/MaraSystems/simple_bank/util"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func NewTestServer(t *testing.T, store db.Store) (*Server, error) {
	config := util.Config{
		AccessSecretKey: util.RandomString(32),
		AccessDuration:  15,
	}
	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server, nil
}

func TestMain(t *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(t.Run())
}
