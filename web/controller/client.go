package controller

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"x-ui/database/model"
	"x-ui/web/service"
	"x-ui/web/session"
)

type ClientController struct {
	InboundController
	clientService service.ClientService
}

func NewClientController(g *gin.RouterGroup) *ClientController {
	a := &ClientController{}
	a.clientRouter(g)
	//a.startTask()
	return a
}

func (a *ClientController) clientRouter(g *gin.RouterGroup) {
	g = g.Group("/client")

	g.GET("/", a.getClients)
	g.POST("/", a.addClient)
	g.GET("/:id", a.getClient)
	g.DELETE("/:id", a.delClient)
	//g.PUT("/:id", a.updateClient)
	//
	//g.POST("/clientIps/:email", a.getClientIps)
	//g.POST("/clearClientIps/:email", a.clearClientIps)
	//g.POST("/resetClientTraffic/:email", a.resetClientTraffic)

}

//func (a *ClientController) startTask() {
//	webServer := global.GetWebServer()
//	c := webServer.GetCron()
//	c.AddFunc("@every 10s", func() {
//		if a.xrayService.IsNeedRestartAndSetFalse() {
//			err := a.xrayService.RestartXray(false)
//			if err != nil {
//				logger.Error("restart xray failed:", err)
//			}
//		}
//	})
//}

func (a *ClientController) getClients(c *gin.Context) {
	user := session.GetLoginUser(c)
	clients, err := a.clientService.GetClients(user.Id)
	if err != nil {
		jsonMsg(c, I18n(c, "pages.clients.toasts.obtain"), err)
		return
	}
	jsonObj(c, clients, nil)
}

func (a *ClientController) getClient(c *gin.Context) {
	clientId := uuid.Must(uuid.FromString(c.Param("id")))
	user := session.GetLoginUser(c)
	clients, err := a.clientService.GetClient(user.Id, clientId)
	if err != nil {
		jsonMsg(c, I18n(c, "pages.clients.toasts.obtain"), err)
		return
	}
	jsonObj(c, clients, nil)
}

//	func (a *InboundController) getInbound(c *gin.Context) {
//		id, err := strconv.Atoi(c.Param("id"))
//		if err != nil {
//			jsonMsg(c, I18n(c, "get"), err)
//			return
//		}
//		inbound, err := a.inboundService.GetInbound(id)
//		if err != nil {
//			jsonMsg(c, I18n(c, "pages.inbounds.toasts.obtain"), err)
//			return
//		}
//		jsonObj(c, inbound, nil)
//	}
func (a *ClientController) addClient(c *gin.Context) {
	var reqBody struct {
		InboundIds []int        `json:"InboundIds" binding:"required"`
		Client     model.Client `json:"Client" binding:"required"`
	}
	var err error
	//client := &model.Client{}

	if err = c.ShouldBind(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		//jsonMsg(c, I18n(c, "pages.clients.addTo"), err)
		return
	}
	var inbounds []*model.Inbound
	test := reqBody.InboundIds
	inbounds, err = a.inboundService.GetInboundsId(test)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

	}

	if len(inbounds) != len(reqBody.InboundIds) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not all inbound IDs exist"})
		return
	}
	for _, inbound := range inbounds {
		reqBody.Client.Inbound = append(reqBody.Client.Inbound, inbound)
	}

	user := session.GetLoginUser(c)
	reqBody.Client.Creator = user.Id
	reqBody.Client.Enable = true
	//Client.Tag = fmt.Sprintf("Client-%v", inbound.Port)
	reqBody.Client, err = a.clientService.AddClient(reqBody.Client)
	jsonMsgObj(c, I18n(c, "pages.clients.addTo"), reqBody.Client, err)
	if err == nil {
		a.xrayService.SetToNeedRestart()
	}
}

func (a *ClientController) delClient(c *gin.Context) {
	clientId := uuid.Must(uuid.FromString(c.Param("id")))
	//user := session.GetLoginUser(c)

	err := a.clientService.DelClient(clientId)
	jsonMsgObj(c, I18n(c, "delete"), clientId, err)
	if err == nil {
		a.xrayService.SetToNeedRestart()
	}
}

//func (a *InboundController) updateInbound(c *gin.Context) {
//	id, err := strconv.Atoi(c.Param("id"))
//	if err != nil {
//		jsonMsg(c, I18n(c, "pages.inbounds.revise"), err)
//		return
//	}
//	inbound := &model.Inbound{
//		Id: id,
//	}
//	err = c.ShouldBind(inbound)
//	if err != nil {
//		jsonMsg(c, I18n(c, "pages.inbounds.revise"), err)
//		return
//	}
//	inbound, err = a.inboundService.UpdateInbound(inbound)
//	jsonMsgObj(c, I18n(c, "pages.inbounds.revise"), inbound, err)
//	if err == nil {
//		a.xrayService.SetToNeedRestart()
//	}
//}
//func (a *InboundController) getClientIps(c *gin.Context) {
//	email := c.Param("email")
//
//	ips, err := a.inboundService.GetInboundClientIps(email)
//	if err != nil {
//		jsonObj(c, "No IP Record", nil)
//		return
//	}
//	jsonObj(c, ips, nil)
//}
//func (a *InboundController) clearClientIps(c *gin.Context) {
//	email := c.Param("email")
//
//	err := a.inboundService.ClearClientIps(email)
//	if err != nil {
//		jsonMsg(c, "修改", err)
//		return
//	}
//	jsonMsg(c, "Log Cleared", nil)
//}
//func (a *InboundController) resetClientTraffic(c *gin.Context) {
//	email := c.Param("email")
//
//	err := a.inboundService.ResetClientTraffic(email)
//	if err != nil {
//		jsonMsg(c, "something worng!", err)
//		return
//	}
//	jsonMsg(c, "traffic reseted", nil)
//}
