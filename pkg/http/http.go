package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/proxy"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

type IPGeoInfo struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	ISP         string  `json:"isp"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Org         string  `json:"org"`
	AS          string  `json:"as"`
}

func Request(url string, data map[string]interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.New("request returned an unexpected status code")
	}

	defer resp.Body.Close()
	return nil
}

func ValidateProxy(proxyAddr, testURL string) (bool, string, error) {
	// 解析代理地址
	proxyURL, err := url.Parse(proxyAddr)
	if err != nil {
		return false, "", fmt.Errorf("代理地址解析失败: %v", err)
	}

	// 获取认证信息
	username := proxyURL.User.Username()
	password, _ := proxyURL.User.Password()

	// 创建自定义 Transport
	var transport *http.Transport
	switch proxyURL.Scheme {
	case "http", "https":
		// HTTP/HTTPS 代理
		transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	case "socks5":
		// SOCKS5 代理处理
		var auth *proxy.Auth
		if username != "" || password != "" {
			auth = &proxy.Auth{
				User:     username,
				Password: password,
			}
		}

		dialer, err := proxy.SOCKS5(
			"tcp",
			net.JoinHostPort(proxyURL.Hostname(), proxyURL.Port()),
			auth,
			proxy.Direct,
		)
		if err != nil {
			return false, "", fmt.Errorf("SOCKS5连接失败: %v", err)
		}

		transport = &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			},
		}
	default:
		return false, "", fmt.Errorf("不支持的代理协议: %s", proxyURL.Scheme)
	}

	// 创建带超时的客户端
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	// 创建请求
	req, _ := http.NewRequest("GET", testURL, nil)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return false, "", fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", fmt.Errorf("响应读取失败: %v", err)
	}

	return resp.StatusCode == http.StatusOK, string(body), nil
}

func GetIPGeoInfo(ip string) (*IPGeoInfo, error) {
	apiURL := fmt.Sprintf("http://ip-api.com/json/%s?fields=status,country,countryCode,region,regionName,city,isp,lat,lon,org,as", ip)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("API请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("响应读取失败: %v", err)
	}

	var geoInfo IPGeoInfo
	if err := json.Unmarshal(body, &geoInfo); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %v", err)
	}

	if geoInfo.Status != "success" {
		return nil, fmt.Errorf("API返回错误: %s", geoInfo.Status)
	}

	return &geoInfo, nil
}
