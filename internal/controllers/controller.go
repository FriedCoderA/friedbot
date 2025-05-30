package controllers

import (
	"errors"

	"friedbot/internal/controllers/receiver"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Controller interface {
	Router(router *gin.RouterGroup)
}

type MainController struct {
	engine      *gin.Engine
	Controllers []Controller
}

func NewMainController() *MainController {
	c := &MainController{
		engine: gin.Default(),
		Controllers: []Controller{
			receiver.NewController(),
		},
	}
	group := c.engine.Group("/napcat/v1")
	c.Route(group)
	return c
}

func (c *MainController) Start() error {
	addr := viper.GetString("server.address")
	if addr == "" {
		return errors.New("server.address is empty")
	}
	return c.engine.Run(addr)
}

func (c *MainController) Route(router *gin.RouterGroup) {
	for _, controller := range c.Controllers {
		controller.Router(router)
	}
}
