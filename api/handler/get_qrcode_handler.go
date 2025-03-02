package handler

import (
	"net/http"
	"zapmeow/api/model"
	"zapmeow/api/response"
	"zapmeow/api/service"
	"zapmeow/pkg/zapmeow"

	"github.com/gin-gonic/gin"
)

type getQrCodeResponse struct {
	QrCode string `json:"qrcode"`
}

type proxyRequest struct {
	ProxyInfo string `json:"proxy"`
}

type getQrCodeHandler struct {
	app              *zapmeow.ZapMeow
	whatsAppService  service.WhatsAppService
	messageService   service.MessageService
	accountService   service.AccountService
	proxyInfoService service.ProxyInfoService
}

func NewGetQrCodeHandler(
	app *zapmeow.ZapMeow,
	whatsAppService service.WhatsAppService,
	messageService service.MessageService,
	accountService service.AccountService,
	proxyInfoService service.ProxyInfoService,
) *getQrCodeHandler {
	return &getQrCodeHandler{
		app:              app,
		whatsAppService:  whatsAppService,
		messageService:   messageService,
		accountService:   accountService,
		proxyInfoService: proxyInfoService,
	}
}

// POST QR Code for WhatsApp Login
//
//	@Summary		Get WhatsApp QR Code
//	@Description	Returns a QR code to initiate WhatsApp login.
//	@Tags			WhatsApp Login
//	@Param			instanceId	path	string	true	"Instance ID"
//	@Param			data		body	proxyRequest	false	"Text message body"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	getQrCodeResponse	"QR Code"
//	@Router			/{instanceId}/qrcode [post]
func (h *getQrCodeHandler) Handler(c *gin.Context) {
	instanceID := c.Param("instanceId")
	var proxyInfo = ""
	var body proxyRequest
	if err := c.ShouldBindJSON(&body); err == nil {
		proxyInfo = body.ProxyInfo
	}
	proxyInfoModel := &model.ProxyInfo{
		User:  instanceID,
		Proxy: proxyInfo,
	}
	proxyErr := h.proxyInfoService.SaveProxyInfo(proxyInfoModel)
	if proxyErr != nil {
		return
	}

	_, err := h.whatsAppService.GetInstance(instanceID, proxyInfo)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

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

	response.Response(c, http.StatusOK, getQrCodeResponse{
		QrCode: account.QrCode,
	})
}
