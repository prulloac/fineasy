package routes

import (
	"github.com/gin-gonic/gin"
)

func Run() *gin.Engine {
	server := gin.Default()
	v1 := server.Group("/v1")
	addPingRoutes(v1)
	addAuthRoutes(v1)
	addSocialRoutes(v1)
	return server
}
