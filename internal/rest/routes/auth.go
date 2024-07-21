package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prulloac/fineasy/internal/middleware"
	"github.com/prulloac/fineasy/internal/rest"
	"github.com/prulloac/fineasy/internal/rest/dto"
)

type AuthController struct {
	authService *middleware.AuthService
}

func NewAuthController(authService *middleware.AuthService) *AuthController {
	instance := &AuthController{}
	instance.authService = authService
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
	var i dto.InternalUserRegisterInput
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rm := rest.GetRequestMeta(c.Request)
	user, err := a.authService.Register(i.Email, i.Password, rm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (a *AuthController) login(c *gin.Context) {
	var i dto.LoginInput
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rm := rest.GetRequestMeta(c.Request)
	out, user, err := a.authService.Login(i.Email, i.Password, rm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Writer.Header().Set("Authorization", rest.GenerateBearerToken(user))
	c.JSON(http.StatusOK, out)
}
