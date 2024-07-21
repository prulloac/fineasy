package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/prulloac/fineasy/internal/db/persistence"
	auth "github.com/prulloac/fineasy/internal/db/repositories/auth"
	core "github.com/prulloac/fineasy/internal/db/repositories/core"
	"github.com/prulloac/fineasy/internal/middleware"
	"github.com/prulloac/fineasy/pkg/logging"
	"github.com/prulloac/fineasy/pkg/validations"
)

type controller interface {
	RegisterEndpoints(rg *gin.RouterGroup)
}

func Server() *gin.Engine {
	server := gin.Default()
	dbconn := persistence.NewPersistence()
	authRepository := auth.NewAuthRepository(dbconn)
	coreRepository := core.NewCoreRepository(dbconn)

	authService := middleware.NewAuthService(authRepository)
	coreService := middleware.NewCoreService(coreRepository)
	authService.NewUserCallbacks(func(u *auth.User) {
		_, err := coreService.CreateUserData(u.ID)
		if err != nil {
			logging.Printf("⚠️ Error creating user data: %s", err)
		}
		g, err := coreService.CreateGroup("Personal", u.ID)
		if err != nil {
			logging.Printf("⚠️ Error creating group: %s", err)
		}
		_, err = coreService.CreateAccount("Personal", "USD", g.ID, u.ID)
		if err != nil {
			logging.Printf("⚠️ Error creating account: %s", err)
		}
		_, err = coreService.CreateUserData(u.ID)
		if err != nil {
			logging.Printf("⚠️ Error creating user data: %s", err)
		}
	})

	controllers := []controller{NewAuthController(authService),
		NewCoreController(coreService)}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("past_time", validations.PastTime)
		v.RegisterValidation("uuid7", validations.UUID7)
		v.RegisterValidation("date", validations.Date)
	}

	v1 := server.Group("/v1")
	v1.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	for _, c := range controllers {
		c.RegisterEndpoints(v1)
	}

	return server
}
