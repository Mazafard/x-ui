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

	g.GET("inbounds/", a.inbounds)
	g.GET("inbounds/:id", a.inbound)
	g.POST("inbounds/", a.addInbound)
	g.DELETE("inbounds/:id", a.delInbound)
	g.PUT("inbounds/:id", a.updateInbound)
	a.inboundController = NewInboundController(g)
	a.clientController = NewClientController(g)

}

func (a *APIController) inbounds(c *gin.Context) {
	a.inboundController.getInbounds(c)
}
func (a *APIController) inbound(c *gin.Context) {
	a.inboundController.getInbound(c)
}
func (a *APIController) addInbound(c *gin.Context) {
	a.inboundController.addInbound(c)
}
func (a *APIController) delInbound(c *gin.Context) {
	a.inboundController.delInbound(c)
}
func (a *APIController) updateInbound(c *gin.Context) {
	a.inboundController.updateInbound(c)
}
