package model

import (
	"gorm.io/gorm"
	"ops-api/utils"
	"time"
)

// DomainServiceProvider 域名服务提供商
type DomainServiceProvider struct {
	Id        uint     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string   `json:"name"`
	AccessKey *string  `json:"access_key"`
	SecretKey *string  `json:"secret_key"`
	Type      uint     `json:"type"`
	Domains   []Domain `gorm:"foreignKey:DomainServiceProviderID"`
}

// BeforeCreate 创建时对敏感信息加密
func (d *DomainServiceProvider) BeforeCreate(tx *gorm.DB) (err error) {
	if d.AccessKey != nil {
		ak, err := utils.Encrypt(*d.AccessKey)
		if err != nil {
			return err
		}
		d.AccessKey = &ak
	}
	if d.SecretKey != nil {
		sk, err := utils.Encrypt(*d.SecretKey)
		if err != nil {
			return err
		}
		d.SecretKey = &sk
	}
	return nil
}

// Domain 域名
type Domain struct {
	Id                      uint                  `json:"id" gorm:"primaryKey;autoIncrement"`
	Name                    string                `json:"name" gorm:"unique"`
	RegistrationAt          *time.Time            `json:"registration_at"`
	ExpirationAt            *time.Time            `json:"expiration_at"`
	DomainServiceProviderID uint                  `json:"domain_service_provider_id"`
	DomainServiceProvider   DomainServiceProvider `json:"domain_service_provider" gorm:"foreignKey:DomainServiceProviderID"`
}
