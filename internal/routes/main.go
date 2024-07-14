package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/prulloac/fineasy/internal/persistence"
)

func Run() *gin.Engine {
	server := gin.Default()
	v1 := server.Group("/v1")
	dbconn := persistence.NewPersistence()
	addPingRoutes(v1)
	NewAuthController(dbconn).RegisterPaths(v1)
	NewSocialController(dbconn).RegisterPaths(v1)
	NewTransactionController(dbconn).RegisterPaths(v1)
	return server
}
