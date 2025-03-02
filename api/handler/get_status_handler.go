package handler

import (
	"net/http"
	"zapmeow/api/response"
	"zapmeow/api/service"
	"zapmeow/pkg/zapmeow"

	"github.com/gin-gonic/gin"
)

type getStatusResponse struct {
	Status string `json:"status"`
	Proxy  string `json:"proxy"`
}

type getStatusHandler struct {
	app              *zapmeow.ZapMeow
	whatsAppService  service.WhatsAppService
	accountService   service.AccountService
	proxyInfoService service.ProxyInfoService
}

func NewGetStatusHandler(
	app *zapmeow.ZapMeow,
	whatsAppService service.WhatsAppService,
	accountService service.AccountService,
	proxyInfoService service.ProxyInfoService,
) *getStatusHandler {
	return &getStatusHandler{
		app:              app,
		whatsAppService:  whatsAppService,
		accountService:   accountService,
		proxyInfoService: proxyInfoService,
	}
}

// Get WhatsApp Instance Status
//
//	@Summary		Get WhatsApp Instance Status
//	@Description	Returns the status of the specified WhatsApp instance.
//	@Tags			WhatsApp Status
//	@Param			instanceId	path	string	true	"Instance ID"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	getStatusResponse	"Status Response"
//	@Router			/{instanceId}/status [get]
func (h *getStatusHandler) Handler(c *gin.Context) {
	instanceID := c.Param("instanceId")
	instance, err := h.whatsAppService.GetInstance(instanceID, "")
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	proxyInfo, _ := h.proxyInfoService.GetProxyInfo(instanceID)

	h.app.Mutex.Lock()
	defer h.app.Mutex.Unlock()
	account, err := h.accountService.GetAccountByInstanceID(instanceID)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if account == nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Account not found")
		return
	}

	var status = account.Status
	if !instance.Client.IsConnected() {
		status = "DISCONNECTED"
	}

	if status == "CONNECTED" && !instance.Client.IsLoggedIn() {
		status = "UNPAIRED"
	}

	response.Response(c, http.StatusOK, getStatusResponse{
		Status: status,
		Proxy:  proxyInfo.Proxy,
	})
}
