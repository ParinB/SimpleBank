package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/parin/simplebank/db/sqlc"
	"os"
	"testing"
)

func TestMain(m *testing.M)  {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func newTestServer(store db.Store) *Server  {
	server  := NewServer(store)
	return server
}