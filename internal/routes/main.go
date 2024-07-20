package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/prulloac/fineasy/internal/persistence"
	"github.com/prulloac/fineasy/pkg/validations"
)

type controller interface {
	Close()
	RegisterEndpoints(rg *gin.RouterGroup)
}

func Server() *gin.Engine {
	server := gin.Default()
	dbconn := persistence.NewPersistence()
	controllers := []controller{NewAuthController(dbconn),
		NewSocialController(dbconn),
		NewTransactionController(dbconn),
		NewPreferencesController(dbconn)}

	registerValidators()

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

func registerValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("past_time", validations.PastTime)
		v.RegisterValidation("uuid7", validations.UUID7)
		v.RegisterValidation("date", validations.Date)
	}
}
