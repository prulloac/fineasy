package routes

import (
	"github.com/gin-gonic/gin"
	p "github.com/prulloac/fineasy/internal/persistence"
	"github.com/prulloac/fineasy/internal/preferences"
)

type PreferencesController struct {
	preferencesService *preferences.Service
}

func NewPreferencesController(persistence *p.Persistence) *PreferencesController {
	return &PreferencesController{preferencesService: preferences.NewService(persistence)}
}

func (c *PreferencesController) Close() {
	c.preferencesService.Close()
}

func (c *PreferencesController) RegisterEndpoints(rg *gin.RouterGroup) {
}
