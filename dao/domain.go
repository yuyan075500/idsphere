package dao

import (
	"ops-api/global"
	"ops-api/model"
)

var Domain domain

type domain struct{}

// ProviderUpdate 更新域名服务商结构体
type ProviderUpdate struct {
	ID        uint    `json:"id" binding:"required"`
	Name      string  `json:"name" binding:"required"`
	AccessKey *string `json:"access_key"`
	SecretKey *string `json:"secret_key"`
	Type      uint    `json:"type" binding:"required"`
}

// AddDomainServiceProvider 新增域名服务商
func (d *domain) AddDomainServiceProvider(data *model.DomainServiceProvider) (provider *model.DomainServiceProvider, err error) {
	if err := global.MySQLClient.Create(&data).Error; err != nil {
		return nil, err
	}
	return provider, nil
}

// DeleteDomainServiceProvider 删除域名服务商
func (d *domain) DeleteDomainServiceProvider(id int) (err error) {
	return global.MySQLClient.Where("id = ?", id).Unscoped().Delete(&model.DomainServiceProvider{}).Error
}

// UpdateDomainServiceProvider 修改域名服务商
func (d *domain) UpdateDomainServiceProvider(updates map[string]interface{}) (*model.DomainServiceProvider, error) {

	// 获取现有的记录
	var provider model.DomainServiceProvider
	if err := global.MySQLClient.First(&provider, updates["id"]).Error; err != nil {
		return nil, err
	}

	// 使用传入的更新字段
	if err := global.MySQLClient.Model(&provider).Updates(updates).Error; err != nil {
		return nil, err
	}

	// 获取更新后的服务商信息并返回
	if err := global.MySQLClient.First(&provider, updates["id"]).Error; err != nil {
		return nil, err
	}
	return &provider, nil
}

// GetDomainServiceProviderList 获取域名服务商列表
func (d *domain) GetDomainServiceProviderList() ([]model.DomainServiceProvider, error) {
	var providers []model.DomainServiceProvider

	// 不返回敏感信息
	if err := global.MySQLClient.Omit("AccessKey", "SecretKey").Find(&providers).Error; err != nil {
		return nil, err
	}

	return providers, nil
}
