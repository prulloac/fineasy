package routes

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/prulloac/fineasy/internal/auth"
	m "github.com/prulloac/fineasy/internal/middleware"
	p "github.com/prulloac/fineasy/internal/persistence"
	"github.com/prulloac/fineasy/internal/social"
	"github.com/prulloac/fineasy/internal/transactions"
	"github.com/prulloac/fineasy/pkg"
)

type AuthController struct {
	authService         *auth.Service
	socialService       *social.Service
	transactionsService *transactions.Service
}

func NewAuthController(persistence *p.Persistence) *AuthController {
	var instance *AuthController
	authService := auth.NewService(persistence)
	socialService := social.NewService(persistence)
	transactionsService := transactions.NewService(persistence)
	authService.NewUserCallbacks(func(u auth.User) {
		g, err := socialService.CreateGroup("Personal", u.ID)
		if err != nil {
			log.Printf("⚠️ Error creating group: %s", err)
		}
		_, err = transactionsService.CreateAccount("Personal", "USD", g.ID, u.ID)
		if err != nil {
			log.Printf("⚠️ Error creating account: %s", err)
		}
	})
	instance = &AuthController{authService: authService, socialService: socialService, transactionsService: transactionsService}
	return instance
}

func (c *AuthController) Close() {
	c.authService.Close()
	c.socialService.Close()
	c.transactionsService.Close()
}

func (c *AuthController) RegisterPaths(rg *gin.RouterGroup) {
	g := rg.Group("/auth")
	g.POST("/register", c.register)
	g.POST("/login", c.login)
	g.GET("/me", m.SecureRequest, c.me)
}

func (a *AuthController) register(c *gin.Context) {
	var i auth.RegisterInput
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := pkg.ValidateStruct(i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rm := pkg.GetRequestMeta(c.Request)
	user, err := a.authService.Register(i.Username, i.Email, i.Password, rm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (a *AuthController) login(c *gin.Context) {
	var i auth.LoginInput
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := pkg.ValidateStruct(i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rm := pkg.GetRequestMeta(c.Request)
	user, err := a.authService.Login(i.Email, i.Password, rm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Writer.Header().Set("Authorization", m.GenerateBearerToken(user))
	c.JSON(http.StatusOK, user)
}

func (a *AuthController) me(c *gin.Context) {
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing token"})
		return
	}
	uid := token.(*jwt.Token).Claims.(jwt.MapClaims)["uid"].(float64)
	user, err := a.authService.Me(uint(uid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}
