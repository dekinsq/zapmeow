package handler

import (
	"github.com/gin-gonic/gin"
	"zapmeow/pkg/zapmeow"
)

type whatsAppInstanceInfo struct {
	InstanceID string `json:"instance_id"`
	//Client     *whatsapp.Instance `json:"client"`
	ProxyInfo string `json:"proxy"`
	Connected bool   `json:"connected"`
	LoggedIn  bool   `json:"logged_in"`
}

type manageInstanceResponse struct {
	Instances []whatsAppInstanceInfo `json:"instances"`
}

type manageInstanceHandler struct {
	app *zapmeow.ZapMeow
}

func NewManageInstanceHandler(app *zapmeow.ZapMeow) *manageInstanceHandler {
	return &manageInstanceHandler{
		app: app,
	}
}

// Manage Instance
//
//	@Summary		Manage Instance
//	@Description	Manage Instance
//	@Tags			Manage Instance
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	manageInstanceResponse	"ManageInstanceResponse"
//	@Router			/manage-instance/all [get]
func (h *manageInstanceHandler) Handler(c *gin.Context) {
	instances := make([]whatsAppInstanceInfo, 0)
	for instanceID, instance := range h.app.GetAllInstances() {
		instances = append(instances, whatsAppInstanceInfo{
			InstanceID: instanceID,
			ProxyInfo:  instance.Client.GetProxyAddress(),
			Connected:  instance.Client.IsConnected(),
			LoggedIn:   instance.Client.IsLoggedIn(),
		})
	}
	c.JSON(200, manageInstanceResponse{
		Instances: instances,
	})
}
