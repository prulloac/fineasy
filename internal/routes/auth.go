package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prulloac/fineasy/internal/auth"
	"github.com/prulloac/fineasy/pkg"
)

func addAuthRoutes(rg *gin.RouterGroup) {
	authGroup := rg.Group("/auth")
	authGroup.POST("/register", register)
	authGroup.POST("/login", login)
}

func register(c *gin.Context) {
	var i auth.RegisterInput
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s := auth.NewService()
	defer s.Close()
	rm := pkg.GetRequestMeta(c.Request)
	user, err := s.Register(i, rm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
	s := auth.NewService()
	defer s.Close()
	rm := pkg.GetRequestMeta(c.Request)
	user, err := s.Login(i, rm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}
