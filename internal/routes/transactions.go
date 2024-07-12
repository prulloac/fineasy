package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	m "github.com/prulloac/fineasy/internal/middleware"
	"github.com/prulloac/fineasy/internal/transactions"
	"github.com/prulloac/fineasy/pkg"
)

func addTransactionsRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/transactions")
	g.POST("accounts", m.SecureRequest, createAccount)
}

func createAccount(c *gin.Context) {
	s := transactions.NewService()
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

	out, err := s.CreateAccount(i.Name, i.GroupID, i.Currency, int(uid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, out)
}
