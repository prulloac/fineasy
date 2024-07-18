package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/prulloac/fineasy/internal/persistence"
	"github.com/prulloac/fineasy/pkg/validations"
)

func Run() *gin.Engine {
	server := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("past_time", validations.PastTime)
		v.RegisterValidation("uuid7", validations.UUID7)
		v.RegisterValidation("date", validations.Date)
	}

	v1 := server.Group("/v1")
	dbconn := persistence.NewPersistence()
	addPingRoutes(v1)
	NewAuthController(dbconn).RegisterPaths(v1)
	NewSocialController(dbconn).RegisterPaths(v1)
	NewTransactionController(dbconn).RegisterPaths(v1)
	return server
}
