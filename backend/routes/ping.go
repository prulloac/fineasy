package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func addPingRoutes(rg *gin.RouterGroup) {
	pingGroup := rg.Group("/ping")

	pingGroup.GET("/", ping)
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
