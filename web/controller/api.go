package controller

import (
	"github.com/gin-gonic/gin"
)

type APIController struct {
	BaseController

	inboundController *InboundController
	clientController  *ClientController
	settingController *SettingController
}

func NewAPIController(g *gin.RouterGroup) *APIController {
	a := &APIController{}
	a.initRouter(g)
	return a
}

func (a *APIController) initRouter(g *gin.RouterGroup) {
	g = g.Group("/xui/API/")
	g.Use(a.checkLogin)
	a.inboundController = NewInboundController(g)
	a.clientController = NewClientController(g)

}
