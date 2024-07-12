package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	m "github.com/prulloac/fineasy/internal/middleware"
	"github.com/prulloac/fineasy/internal/social"
	"github.com/prulloac/fineasy/pkg"
)

func addSocialRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/social")
	g.POST("/friends/requests", m.SecureRequest, addFriend)
	g.GET("/friends", m.SecureRequest, getFriends)
	g.GET("/friends/requests", m.SecureRequest, getFriendRequests)
	g.PATCH("/friends/requests", m.SecureRequest, updateFriendRequest)
	g.DELETE("/friends", m.SecureRequest, deleteFriend)
	g.POST("/groups", m.SecureRequest, createGroup)
	g.GET("/groups", m.SecureRequest, getUserGroups)
	g.PATCH("/groups", m.SecureRequest, updateGroup)
	g.POST("/groups/membership", m.SecureRequest, updateUserGroup)
}

func addFriend(c *gin.Context) {
	s := social.NewService()
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client user not found"})
		return
	}
	uid := token.(*jwt.Token).Claims.(jwt.MapClaims)["uid"].(float64)

	var i social.AddFriendInput
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := pkg.ValidateStruct(i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if i.FriendID == int(uid) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot add yourself as friend"})
		return
	}
	if i.UserID != int(uid) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	out, err := s.AddFriend(i.FriendID, int(uid))
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
	uid := token.(*jwt.Token).Claims.(jwt.MapClaims)["uid"].(float64)

	out, err := s.GetFriends(int(uid))
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
	uid := token.(*jwt.Token).Claims.(jwt.MapClaims)["uid"].(float64)

	out, err := s.GetFriendRequests(int(uid))
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
	uid := token.(*jwt.Token).Claims.(jwt.MapClaims)["uid"].(float64)

	var i social.UpdateFriendRequestInput
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := pkg.ValidateStruct(i); err != nil || i.FriendID == int(uid) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if i.UserID != int(uid) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	if i.Status != "Accepted" && i.Status != "Rejected" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}

	out, err := s.UpdateFriendRequest(i.FriendID, i.Status, int(uid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func deleteFriend(c *gin.Context) {
	s := social.NewService()
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client user not found"})
		return
	}
	uid := token.(*jwt.Token).Claims.(jwt.MapClaims)["uid"].(float64)

	var i social.DeleteFriendInput
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := pkg.ValidateStruct(i); err != nil || i.FriendID == int(uid) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if i.UserID != int(uid) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	out, err := s.DeleteFriend(i.FriendID, int(uid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, out)
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
	if err := pkg.ValidateStruct(i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid := token.(*jwt.Token).Claims.(jwt.MapClaims)["uid"].(float64)
	out, err := s.CreateGroup(i.Name, int(uid))
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

	c.JSON(http.StatusOK, out)
}
