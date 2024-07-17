package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	m "github.com/prulloac/fineasy/internal/middleware"
	p "github.com/prulloac/fineasy/internal/persistence"
	"github.com/prulloac/fineasy/internal/social"
)

type SocialController struct {
	socialService *social.Service
}

func NewSocialController(persistence *p.Persistence) *SocialController {
	return &SocialController{socialService: social.NewService(persistence)}
}

func (c *SocialController) Close() {
	c.socialService.Close()
}

func (c *SocialController) RegisterPaths(rg *gin.RouterGroup) {
	p := rg.Group("/social", m.SecureRequest)
	// friends
	f := p.Group("/friends")
	f.GET("", c.getFriendships)
	f.GET(":id", c.getFriendhip)
	f.DELETE(":id", c.deleteFriendship)
	f.POST("/requests", c.addFriendship)
	f.GET("/requests", c.getPendingFriendships)
	f.PATCH("/requests/:id", c.acceptFriendship)
	// groups
	g := p.Group("/groups")
	g.POST("", c.createGroup)
	g.GET("", c.getUserGroups)
	g.GET(":id", c.getUserGroup)
	g.PATCH(":id", c.updateGroup)
	g.POST("/invite", c.inviteUserGroup)
	g.PATCH("/membership", c.updateUserGroup)
}

func (s *SocialController) addFriendship(c *gin.Context) {
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
	if i.FriendID == uint(uid) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot add yourself as friend"})
		return
	}
	if i.UserID != uint(uid) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	out, err := s.socialService.AddFriendship(i.FriendID, uint(uid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, out)
}

func (s *SocialController) getFriendships(c *gin.Context) {
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client user not found"})
		return
	}
	uid := token.(*jwt.Token).Claims.(jwt.MapClaims)["uid"].(float64)

	out, err := s.socialService.GetFriendships(uint(uid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *SocialController) getFriendhip(c *gin.Context) {
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client user not found"})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid := token.(*jwt.Token).Claims.(jwt.MapClaims)["uid"].(float64)

	out, err := s.socialService.GetFriendship(uint(id), uint(uid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *SocialController) getPendingFriendships(c *gin.Context) {
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client user not found"})
		return
	}
	uid := token.(*jwt.Token).Claims.(jwt.MapClaims)["uid"].(float64)

	out, err := s.socialService.GetPendingFriendships(uint(uid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *SocialController) acceptFriendship(c *gin.Context) {
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
	if i.Status != "Accepted" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}
	fid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := s.socialService.AcceptFriendship(i.Status, uint(fid), uint(uid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *SocialController) deleteFriendship(c *gin.Context) {
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
	if i.FriendID == uint(uid) {
		c.JSON(http.StatusForbidden, gin.H{"error": "cannot delete yourself as friend"})
		return
	}
	if i.UserID != uint(uid) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	out, err := s.socialService.RejectFriendship(i.FriendID, uint(uid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, out)
}

func (s *SocialController) createGroup(c *gin.Context) {
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
	uid := token.(*jwt.Token).Claims.(jwt.MapClaims)["uid"].(float64)
	out, err := s.socialService.CreateGroup(i.Name, uint(uid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, out)
}

func (s *SocialController) getUserGroups(c *gin.Context) {
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client user not found"})
		return
	}
	uid := token.(*jwt.Token).Claims.(jwt.MapClaims)["uid"].(float64)

	out, err := s.socialService.GetUserGroups(uint(uid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *SocialController) getUserGroup(c *gin.Context) {
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client user not found"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uid := token.(*jwt.Token).Claims.(jwt.MapClaims)["uid"].(float64)
	out, err := s.socialService.GetGroupByID(uint(id), uint(uid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *SocialController) updateGroup(c *gin.Context) {
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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid := token.(*jwt.Token).Claims.(jwt.MapClaims)["uid"].(float64)

	out, err := s.socialService.UpdateGroup(i.Name, uint(id), uint(uid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *SocialController) inviteUserGroup(c *gin.Context) {
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
	if i.Status != "Invited" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}
	uid := token.(*jwt.Token).Claims.(jwt.MapClaims)["uid"].(float64)
	if i.UserID == uint(uid) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot invite yourself"})
		return
	}

	out, err := s.socialService.InviteUserGroup(i.GroupID, i.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, out)
}

func (s *SocialController) updateUserGroup(c *gin.Context) {
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
	if i.Status == "Invited" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}
	uid := token.(*jwt.Token).Claims.(jwt.MapClaims)["uid"].(float64)
	if i.UserID != uint(uid) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	out, err := s.socialService.UpdateUserGroup(i.Status, i.GroupID, uint(uid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}
