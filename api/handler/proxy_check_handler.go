package handler

import (
	"github.com/gin-gonic/gin"
	"zapmeow/api/response"
	"zapmeow/pkg/http"
)

type proxyCheckResponse struct {
	Success bool           `json:"success"`
	Ip      string         `json:"ip"`
	IpInfo  http.IPGeoInfo `json:"ipinfo"`
}

type proxyCheckRequest struct {
	ProxyInfo string `json:"proxy"`
}

type proxyCheckHandler struct {
}

func NewProxyCheckHandler() *proxyCheckHandler {
	return &proxyCheckHandler{}
}

// Proxy Check
//
//	@Summary		Proxy Check
//	@Description	Proxy Check
//	@Tags			Proxy Check
//	@Param			data		body	proxyCheckRequest	true	"Text message body"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	proxyCheckResponse	"ProxyCheckResponse"
//	@Router			/proxy-check [post]
func (h *proxyCheckHandler) Handler(c *gin.Context) {
	var body proxyCheckRequest
	var proxyInfo = ""
	if err := c.ShouldBindJSON(&body); err == nil {
		proxyInfo = body.ProxyInfo
	}
	success, ip, err := http.ValidateProxy(proxyInfo, "https://api.ipify.org")
	if err == nil {
		ipinfo, _ := http.GetIPGeoInfo(ip)
		response.Response(c, 200, proxyCheckResponse{
			Success: success,
			Ip:      ip,
			IpInfo:  *ipinfo,
		})
	} else {
		response.ErrorResponse(c, 500, err.Error())
	}
}
