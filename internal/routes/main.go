package api

import (
	"github.com/gin-gonic/gin"
)

func Run() {
	server := gin.Default()
	v1 := server.Group("/v1")
	addPingRoutes(v1)
	server.Run()
}
