package repository

import (
	"gorm.io/gorm"
	"zapmeow/api/model"
	"zapmeow/pkg/database"
)

type ProxyInfoRepository interface {
	InsertProxyInfo(proxyInfo *model.ProxyInfo) error
	UpdateProxyInfo(user string, data map[string]interface{}) error
	GetProxyInfo(user string) (*model.ProxyInfo, error)
	DeleteProxyInfo(user string) error
}

type proxyInfoRepository struct {
	database database.Database
}

func NewProxyInfoRepository(database database.Database) *proxyInfoRepository {
	return &proxyInfoRepository{database: database}
}

func (repo *proxyInfoRepository) InsertProxyInfo(proxyInfo *model.ProxyInfo) error {
	return repo.database.Client().Create(proxyInfo).Error
}

func (repo *proxyInfoRepository) UpdateProxyInfo(user string, data map[string]interface{}) error {
	var proxyInfo model.ProxyInfo
	if result := repo.database.Client().Where("user=?", user).First(&proxyInfo); result.Error != nil {
		return result.Error
	}
	if err := repo.database.Client().Model(&proxyInfo).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func (repo *proxyInfoRepository) GetProxyInfo(user string) (*model.ProxyInfo, error) {
	var proxyInfo model.ProxyInfo
	result := repo.database.Client().Where("user=?", user).First(&proxyInfo)
	if result.Error != nil {
		if result.Error != gorm.ErrRecordNotFound {
			return nil, result.Error
		}
		return nil, nil
	}
	return &proxyInfo, nil
}

func (repo *proxyInfoRepository) DeleteProxyInfo(user string) error {
	if result := repo.database.Client().Where("user=?", user).Unscoped().Delete(&model.ProxyInfo{}); result.Error != nil {
		return result.Error
	}
	return nil
}
