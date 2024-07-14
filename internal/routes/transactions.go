package routes

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	m "github.com/prulloac/fineasy/internal/middleware"
	p "github.com/prulloac/fineasy/internal/persistence"
	"github.com/prulloac/fineasy/internal/transactions"
	"github.com/prulloac/fineasy/pkg"
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
	g := rg.Group("/transactions")
	g.POST("accounts", m.SecureRequest, c.createAccount)
	g.GET("accounts", m.SecureRequest, c.getAccounts)
	g.GET("accounts/:id", m.SecureRequest, c.getAccount)
	g.PATCH("accounts/:id", m.SecureRequest, c.updateAccount)
	// g.DELETE("accounts/:id", m.SecureRequest, deleteAccount)
	// g.POST("categories", m.SecureRequest, createCategory)
	// g.GET("categories", m.SecureRequest, getCategories)
	// g.GET("categories/:id", m.SecureRequest, getCategory)
	// g.PATCH("categories/:id", m.SecureRequest, updateCategory)
	// g.DELETE("categories/:id", m.SecureRequest, deleteCategory)
	// g.POST("transactions", m.SecureRequest, createTransaction)
	// g.GET("transactions", m.SecureRequest, getTransactions)
	// g.GET("transactions/:id", m.SecureRequest, getTransaction)
	// g.PATCH("transactions/:id", m.SecureRequest, updateTransaction)
	// g.DELETE("transactions/:id", m.SecureRequest, deleteTransaction)
	// g.POST("budgets", m.SecureRequest, createBudget)
	// g.GET("budgets", m.SecureRequest, getBudgets)
	// g.GET("budgets/:id", m.SecureRequest, getBudget)
	// g.PATCH("budgets/:id", m.SecureRequest, updateBudget)
	// g.DELETE("budgets/:id", m.SecureRequest, deleteBudget)
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
	if err := pkg.ValidateStruct(i); err != nil {
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
	if err := pkg.ValidateStruct(i); err != nil {
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
