package routes

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/prulloac/fineasy/internal/middleware"
	"github.com/prulloac/fineasy/internal/rest"
	"github.com/prulloac/fineasy/internal/rest/dto"
)

type CoreController struct {
	coreService *middleware.CoreService
}

func NewCoreController(coreService *middleware.CoreService) *CoreController {
	instance := &CoreController{}
	instance.coreService = coreService
	return instance
}

func (c *CoreController) Close() {
	c.coreService.Close()
}

func (c *CoreController) RegisterEndpoints(rg *gin.RouterGroup) {
	// friends
	f := rg.Group("/friends")
	f.Use(rest.CaptureTokenFromHeader)
	f.GET("", c.getFriendships)
	f.GET(":id", c.getFriendhip)
	f.DELETE(":id", c.deleteFriendship)
	f.POST("/requests", c.addFriendship)
	f.GET("/requests", c.getPendingFriendships)
	f.PATCH("/requests/:id", c.acceptFriendship)
	// groups
	g := rg.Group("/groups")
	g.Use(rest.CaptureTokenFromHeader)
	g.POST("", c.createGroup)
	g.GET("", c.getUserGroups)
	g.GET(":id", c.getUserGroup)
	g.PATCH(":id", c.updateGroup)
	g.POST("/invite", c.inviteUserGroup)
	g.PATCH("/membership", c.updateUserGroup)
	// accounts
	a := rg.Group("/accounts")
	a.Use(rest.CaptureTokenFromHeader)
	a.POST("", c.createAccount)
	a.GET("", c.getAccounts)
	a.GET("/:id", c.getAccount)
	a.PATCH("/:id", c.updateAccount)
	// g.DELETE("accounts/:id", c.deleteAccount)
	b := rg.Group("/budgets")
	b.Use(rest.CaptureTokenFromHeader)
	b.POST("", c.createBudget)
	b.GET("", c.getBudgets)
	// b.GET("/:id", c.getBudget)
	// b.PATCH("/:id", c.updateBudget)
	// b.DELETE("/:id", c.deleteBudget)
	t := rg.Group("/transactions")
	t.Use(rest.CaptureTokenFromHeader)
	t.POST("", c.createTransaction)
	// t.GET("transactions", c.getTransactions)
	// t.GET("transactions/:id", c.getTransaction)
	// t.PATCH("transactions/:id", c.updateTransaction)
	// t.DELETE("transactions/:id", c.deleteTransaction)
}

func (s *CoreController) addFriendship(c *gin.Context) {
	token := rest.GetClientToken(c)
	uid := token.Claims.(jwt.MapClaims)["uid"].(float64)

	var i dto.AddFriendInput
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
	out, err := s.coreService.AddFriendship(i.FriendID, uint(uid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, out)
}

func (s *CoreController) getFriendships(c *gin.Context) {
	token := rest.GetClientToken(c)
	uid := token.Claims.(jwt.MapClaims)["uid"].(float64)

	out, err := s.coreService.GetFriendships(uint(uid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *CoreController) getFriendhip(c *gin.Context) {
	token := rest.GetClientToken(c)
	uid := token.Claims.(jwt.MapClaims)["uid"].(float64)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := s.coreService.GetFriendship(uint(id), uint(uid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *CoreController) getPendingFriendships(c *gin.Context) {
	token := rest.GetClientToken(c)
	uid := token.Claims.(jwt.MapClaims)["uid"].(float64)

	out, err := s.coreService.GetPendingFriendships(uint(uid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *CoreController) acceptFriendship(c *gin.Context) {
	token := rest.GetClientToken(c)
	uid := token.Claims.(jwt.MapClaims)["uid"].(float64)

	var i dto.UpdateFriendRequestInput
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

	out, err := s.coreService.AcceptFriendship(i.Status, uint(fid), uint(uid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *CoreController) deleteFriendship(c *gin.Context) {
	token := rest.GetClientToken(c)
	uid := token.Claims.(jwt.MapClaims)["uid"].(float64)

	var i dto.DeleteFriendInput
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

	out, err := s.coreService.RejectFriendship(i.FriendID, uint(uid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, out)
}

func (s *CoreController) createGroup(c *gin.Context) {
	token := rest.GetClientToken(c)
	uid := token.Claims.(jwt.MapClaims)["uid"].(float64)

	var i dto.CreateGroupInput
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out, err := s.coreService.CreateGroup(i.Name, uint(uid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, out)
}

func (s *CoreController) getUserGroups(c *gin.Context) {
	token := rest.GetClientToken(c)
	uid := token.Claims.(jwt.MapClaims)["uid"].(float64)

	out, err := s.coreService.GetUserGroups(uint(uid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *CoreController) getUserGroup(c *gin.Context) {
	token := rest.GetClientToken(c)
	uid := token.Claims.(jwt.MapClaims)["uid"].(float64)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := s.coreService.GetGroupByID(uint(id), uint(uid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *CoreController) updateGroup(c *gin.Context) {
	token := rest.GetClientToken(c)
	uid := token.Claims.(jwt.MapClaims)["uid"].(float64)

	var i dto.UpdateGroupInput
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := s.coreService.UpdateGroup(i.Name, uint(id), uint(uid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *CoreController) inviteUserGroup(c *gin.Context) {
	token := rest.GetClientToken(c)
	uid := token.Claims.(jwt.MapClaims)["uid"].(float64)

	var i dto.JoinGroupInput
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if i.Status != "Invited" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}
	if i.UserID == uint(uid) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot invite yourself"})
		return
	}

	out, err := s.coreService.InviteUserGroup(i.GroupID, i.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, out)
}

func (s *CoreController) updateUserGroup(c *gin.Context) {
	token := rest.GetClientToken(c)
	uid := token.Claims.(jwt.MapClaims)["uid"].(float64)

	var i dto.JoinGroupInput
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if i.Status == "Invited" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}
	if i.UserID != uint(uid) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	out, err := s.coreService.UpdateUserGroup(i.Status, i.GroupID, uint(uid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (t *CoreController) createAccount(c *gin.Context) {
	token := rest.GetClientToken(c)
	uid := token.Claims.(jwt.MapClaims)["uid"].(float64)
	var i dto.CreateAccountInput
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := t.coreService.CreateAccount(i.Name, i.Currency, i.GroupID, uint(uid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, out)
}

func (t *CoreController) getAccounts(c *gin.Context) {
	token := rest.GetClientToken(c)
	uid := token.Claims.(jwt.MapClaims)["uid"].(float64)

	out, err := t.coreService.GetAccounts(uint(uid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (t *CoreController) getAccount(c *gin.Context) {
	token := rest.GetClientToken(c)
	uid := token.Claims.(jwt.MapClaims)["uid"].(float64)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := t.coreService.GetAccountByID(uint(id), uint(uid))
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, out)
}

func (t *CoreController) updateAccount(c *gin.Context) {
	token := rest.GetClientToken(c)
	uid := token.Claims.(jwt.MapClaims)["uid"].(float64)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var i dto.UpdateAccountInput
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	balance, err := strconv.ParseFloat(i.Balance, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := t.coreService.UpdateAccount(i.Name, i.Currency, balance, uint(id), uint(uid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (t *CoreController) createBudget(c *gin.Context) {
	token := rest.GetClientToken(c)
	uid := token.Claims.(jwt.MapClaims)["uid"].(float64)

	var i dto.CreateBudgetInput
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	amount, err := strconv.ParseFloat(i.Amount, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	startDate, err := time.Parse(time.DateOnly, i.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	endDate, err := time.Parse(time.DateOnly, i.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := t.coreService.CreateBudget(i.Name, i.Currency, amount, startDate, endDate, i.AccountID, uint(uid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, out)
}

func (t *CoreController) getBudgets(c *gin.Context) {
	token := rest.GetClientToken(c)
	uid := token.Claims.(jwt.MapClaims)["uid"].(float64)

	out, err := t.coreService.GetBudgets(uint(uid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (t *CoreController) createTransaction(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
