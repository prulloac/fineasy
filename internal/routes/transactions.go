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

func (c *TransactionController) RegisterPaths(rg *gin.RouterGroup) {
	g := rg.Group("/transactions", m.SecureRequest)
	ac := g.Group("/accounts")
	ac.POST("", c.createAccount)
	ac.GET("", c.getAccounts)
	ac.GET("/:id", c.getAccount)
	ac.PATCH("/:id", c.updateAccount)
	// g.DELETE("accounts/:id", c.deleteAccount)
	b := g.Group("/budgets")
	b.POST("", c.createBudget)
	// g.GET("", c.getBudgets)
	// g.GET("/:id", c.getBudget)
	// g.PATCH("/:id", c.updateBudget)
	// g.DELETE("/:id", c.deleteBudget)
	tx := g.Group("/transactions")
	tx.POST("", c.createTransaction)
	// g.GET("transactions", c.getTransactions)
	// g.GET("transactions/:id", c.getTransaction)
	// g.PATCH("transactions/:id", c.updateTransaction)
	// g.DELETE("transactions/:id", c.deleteTransaction)
}

func (t *TransactionController) createAccount(c *gin.Context) {
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client user not found"})
		return
	}

	var i transactions.CreateAccountInput
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid := token.(*jwt.Token).Claims.(jwt.MapClaims)["uid"].(float64)

	out, err := t.transactionService.CreateAccount(i.Name, i.Currency, i.GroupID, uint(uid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, out)
}

func (t *TransactionController) getAccounts(c *gin.Context) {
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client user not found"})
		return
	}

	uid := uint(token.(*jwt.Token).Claims.(jwt.MapClaims)["uid"].(float64))

	out, err := t.transactionService.GetAccounts(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (t *TransactionController) getAccount(c *gin.Context) {
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
	uid := int(token.(*jwt.Token).Claims.(jwt.MapClaims)["uid"].(float64))

	out, err := t.transactionService.GetAccountByID(id, uid)
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
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client user not found"})
		return
	}

	uid := int(token.(*jwt.Token).Claims.(jwt.MapClaims)["uid"].(float64))
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

	out, err := t.transactionService.UpdateAccount(i.Name, i.Currency, balance, id, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (t *TransactionController) createBudget(c *gin.Context) {
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client user not found"})
		return
	}

	var i transactions.CreateBudgetInput
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid := token.(*jwt.Token).Claims.(jwt.MapClaims)["uid"].(float64)
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

func (t *TransactionController) createTransaction(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
