package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	m "github.com/prulloac/fineasy/internal/middleware"
	"github.com/prulloac/fineasy/internal/social"
)

func addSocialRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/social")
	g.POST("/friends", m.SecureRequest, addFriend)
	g.GET("/friends", m.SecureRequest, getFriends)
	g.PATCH("/friends", m.SecureRequest, updateFriend)
}

func addFriend(c *gin.Context) {
	s := social.NewService()
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client user not found"})
		return
	}

	var i social.AddFriendInput
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := s.AddFriend(i, token.(*jwt.Token))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, out)
}

func getFriends(c *gin.Context) {
}

func updateFriend(c *gin.Context) {
}
