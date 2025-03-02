package service

import (
	"zapmeow/api/model"
	"zapmeow/api/repository"
)

type ProxyInfoService interface {
	SaveProxyInfo(proxyInfo *model.ProxyInfo) error
	InsertProxyInfo(proxyInfo *model.ProxyInfo) error
	UpdateProxyInfo(user string, data map[string]interface{}) error
	GetProxyInfo(user string) (*model.ProxyInfo, error)
	DeleteProxyInfo(user string) error
}

type proxyInfoService struct {
	proxyInfoRepo repository.ProxyInfoRepository
}

func NewProxyInfoService(proxyInfoRepo repository.ProxyInfoRepository) *proxyInfoService {
	return &proxyInfoService{
		proxyInfoRepo: proxyInfoRepo,
	}
}

func (p *proxyInfoService) SaveProxyInfo(proxyInfo *model.ProxyInfo) error {
	if proxyInfo != nil && proxyInfo.User != "" {
		existProxyInfo, _ := p.GetProxyInfo(proxyInfo.User)
		if existProxyInfo == nil {
			err := p.proxyInfoRepo.InsertProxyInfo(proxyInfo)
			if err != nil {
				return err
			}
		} else {
			err := p.proxyInfoRepo.UpdateProxyInfo(proxyInfo.User, map[string]interface{}{
				"proxy": proxyInfo.Proxy,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *proxyInfoService) InsertProxyInfo(proxyInfo *model.ProxyInfo) error {
	return p.proxyInfoRepo.InsertProxyInfo(proxyInfo)
}

func (p *proxyInfoService) UpdateProxyInfo(user string, data map[string]interface{}) error {
	return p.UpdateProxyInfo(user, data)
}

func (p *proxyInfoService) GetProxyInfo(user string) (*model.ProxyInfo, error) {
	return p.proxyInfoRepo.GetProxyInfo(user)
}

func (p *proxyInfoService) DeleteProxyInfo(user string) error {
	return p.proxyInfoRepo.DeleteProxyInfo(user)
}
