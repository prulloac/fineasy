package routes

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	m "github.com/prulloac/fineasy/internal/middleware"
	p "github.com/prulloac/fineasy/internal/persistence"
	"github.com/prulloac/fineasy/internal/transactions"
)

type TransactionController struct {
	transactionService *transactions.Service
}

func NewTransactionController(persistence *p.Persistence) *TransactionController {
	return &TransactionController{transactionService: transactions.NewService(persistence)}
}

func (c *TransactionController) Close() {
	c.transactionService.Close()
}

func (c *TransactionController) RegisterEndpoints(rg *gin.RouterGroup) {
	a := rg.Group("/accounts")
	a.Use(m.CaptureTokenFromHeader)
	a.POST("", c.createAccount)
	a.GET("", c.getAccounts)
	a.GET("/:id", c.getAccount)
	a.PATCH("/:id", c.updateAccount)
	// g.DELETE("accounts/:id", c.deleteAccount)
	b := rg.Group("/budgets")
	b.Use(m.CaptureTokenFromHeader)
	b.POST("", c.createBudget)
	b.GET("", c.getBudgets)
	// b.GET("/:id", c.getBudget)
	// b.PATCH("/:id", c.updateBudget)
	// b.DELETE("/:id", c.deleteBudget)
	t := rg.Group("/transactions")
	t.Use(m.CaptureTokenFromHeader)
	t.POST("", c.createTransaction)
	// t.GET("transactions", c.getTransactions)
	// t.GET("transactions/:id", c.getTransaction)
	// t.PATCH("transactions/:id", c.updateTransaction)
	// t.DELETE("transactions/:id", c.deleteTransaction)
}

func (t *TransactionController) createAccount(c *gin.Context) {
	token := m.GetClientToken(c)
	uid := token.Claims.(jwt.MapClaims)["uid"].(float64)
	var i transactions.CreateAccountInput
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := t.transactionService.CreateAccount(i.Name, i.Currency, i.GroupID, uint(uid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, out)
}

func (t *TransactionController) getAccounts(c *gin.Context) {
	token := m.GetClientToken(c)
	uid := token.Claims.(jwt.MapClaims)["uid"].(float64)

	out, err := t.transactionService.GetAccounts(uint(uid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (t *TransactionController) getAccount(c *gin.Context) {
	token := m.GetClientToken(c)
	uid := token.Claims.(jwt.MapClaims)["uid"].(float64)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := t.transactionService.GetAccountByID(uint(id), uint(uid))
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

func (t *TransactionController) updateAccount(c *gin.Context) {
	token := m.GetClientToken(c)
	uid := token.Claims.(jwt.MapClaims)["uid"].(float64)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var i transactions.UpdateAccountInput
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	balance, err := strconv.ParseFloat(i.Balance, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := t.transactionService.UpdateAccount(i.Name, i.Currency, balance, uint(id), uint(uid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (t *TransactionController) createBudget(c *gin.Context) {
	token := m.GetClientToken(c)
	uid := token.Claims.(jwt.MapClaims)["uid"].(float64)

	var i transactions.CreateBudgetInput
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

	out, err := t.transactionService.CreateBudget(i.Name, i.Currency, amount, startDate, endDate, i.AccountID, uint(uid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, out)
}

func (t *TransactionController) getBudgets(c *gin.Context) {
	token := m.GetClientToken(c)
	uid := token.Claims.(jwt.MapClaims)["uid"].(float64)

	out, err := t.transactionService.GetBudgets(uint(uid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (t *TransactionController) createTransaction(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
