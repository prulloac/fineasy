package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prulloac/fineasy/internal/auth"
	m "github.com/prulloac/fineasy/internal/middleware"
	p "github.com/prulloac/fineasy/internal/persistence"
	"github.com/prulloac/fineasy/internal/preferences"
	"github.com/prulloac/fineasy/internal/social"
	"github.com/prulloac/fineasy/internal/transactions"
	"github.com/prulloac/fineasy/pkg"
	"github.com/prulloac/fineasy/pkg/logging"
)

type AuthController struct {
	authService *auth.Service
}

func NewAuthController(persistence *p.Persistence) *AuthController {
	var instance *AuthController
	authService := auth.NewService(persistence)
	socialService := social.NewService(persistence)
	transactionsService := transactions.NewService(persistence)
	preferencesService := preferences.NewService(persistence)
	authService.NewUserCallbacks(func(u *auth.User) {
		g, err := socialService.CreateGroup("Personal", u.ID)
		if err != nil {
			logging.Printf("⚠️ Error creating group: %s", err)
		}
		_, err = transactionsService.CreateAccount("Personal", "USD", g.ID, u.ID)
		if err != nil {
			logging.Printf("⚠️ Error creating account: %s", err)
		}
		_, err = preferencesService.CreateUserData(u.ID)
		if err != nil {
			logging.Printf("⚠️ Error creating user data: %s", err)
		}
	})
	instance = &AuthController{authService: authService}
	return instance
}

func (c *AuthController) Close() {
	c.authService.Close()
}

func (c *AuthController) RegisterEndpoints(rg *gin.RouterGroup) {
	g := rg.Group("/auth")
	g.POST("/register", c.register)
	g.POST("/login", c.login)
}

func (a *AuthController) register(c *gin.Context) {
	var i auth.InternalUserRegisterInput
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rm := pkg.GetRequestMeta(c.Request)
	user, err := a.authService.Register(i.Email, i.Password, rm)
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
	rm := pkg.GetRequestMeta(c.Request)
	out, user, err := a.authService.Login(i.Email, i.Password, rm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Writer.Header().Set("Authorization", m.GenerateBearerToken(user))
	c.JSON(http.StatusOK, out)
}
