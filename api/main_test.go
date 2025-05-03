package api

import (
	"testing"

	db "github.com/Sidsha242/simple_bank/db/sqlc"
	"github.com/Sidsha242/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		JWTSecret:     util.RandomString(32),
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}