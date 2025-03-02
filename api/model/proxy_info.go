package model

import "gorm.io/gorm"

type ProxyInfo struct {
	gorm.Model
	User  string
	Proxy string
}
