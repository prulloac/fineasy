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
	g.POST("/friends/requests", m.SecureRequest, addFriend)
	g.GET("/friends", m.SecureRequest, getFriends)
	g.GET("/friends/requests", m.SecureRequest, getFriendRequests)
	g.PATCH("/friends/requests", m.SecureRequest, updateFriendRequest)
	// g.DELETE("/friends", m.SecureRequest, deleteFriend)
	g.POST("/groups", m.SecureRequest, createGroup)
	g.GET("/groups", m.SecureRequest, getUserGroups)
	g.PATCH("/groups", m.SecureRequest, updateGroup)
	g.POST("/groups/membership", m.SecureRequest, updateUserGroup)
	// g.DELETE("/groups", m.SecureRequest, dropGroup)
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
	s := social.NewService()
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client user not found"})
		return
	}

	out, err := s.GetFriends(token.(*jwt.Token))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func getFriendRequests(c *gin.Context) {
	s := social.NewService()
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client user not found"})
		return
	}

	out, err := s.GetFriendRequests(token.(*jwt.Token))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func updateFriendRequest(c *gin.Context) {
	s := social.NewService()
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client user not found"})
		return
	}

	var i social.UpdateFriendRequestInput
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := s.UpdateFriendRequest(i, token.(*jwt.Token))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func createGroup(c *gin.Context) {
	s := social.NewService()
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client user not found"})
		return
	}

	var i social.CreateGroupInput
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := s.CreateGroup(i, token.(*jwt.Token))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, out)
}

func getUserGroups(c *gin.Context) {
	s := social.NewService()
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client user not found"})
		return
	}

	out, err := s.GetUserGroups(token.(*jwt.Token))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func updateGroup(c *gin.Context) {
	s := social.NewService()
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client user not found"})
		return
	}

	var i social.UpdateGroupInput
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := s.UpdateGroup(i, token.(*jwt.Token))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func updateUserGroup(c *gin.Context) {
	s := social.NewService()
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client user not found"})
		return
	}

	var i social.JoinGroupInput
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := s.UpdateUserGroup(i, token.(*jwt.Token))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, out)
}
