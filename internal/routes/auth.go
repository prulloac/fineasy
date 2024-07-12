package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/prulloac/fineasy/internal/auth"
	"github.com/prulloac/fineasy/internal/middleware"
	m "github.com/prulloac/fineasy/internal/middleware"
	"github.com/prulloac/fineasy/internal/social"
	"github.com/prulloac/fineasy/pkg"
)

func addAuthRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/auth")
	g.POST("/register", register)
	g.POST("/login", login)
	g.GET("/me", m.SecureRequest, me)
}

func register(c *gin.Context) {
	var i auth.RegisterInput
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := pkg.ValidateStruct(i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s := auth.NewService()
	rm := pkg.GetRequestMeta(c.Request)
	user, err := s.Register(i.Username, i.Email, i.Password, rm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err = social.NewService().CreateGroup("Personal", user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func login(c *gin.Context) {
	var i auth.LoginInput
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := pkg.ValidateStruct(i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s := auth.NewService()
	rm := pkg.GetRequestMeta(c.Request)
	user, err := s.Login(i.Email, i.Password, rm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Writer.Header().Set("Authorization", middleware.GenerateBearerToken(user))
	c.JSON(http.StatusOK, user)
}

func me(c *gin.Context) {
	s := auth.NewService()
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing token"})
		return
	}
	uhash := token.(*jwt.Token).Claims.(jwt.MapClaims)["sub"].(string)
	user, err := s.Me(uhash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}
