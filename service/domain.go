package service

import (
	"ops-api/dao"
	"ops-api/model"
	"ops-api/utils"
)

var Domain domain

type domain struct{}

// DomainServiceProviderCreate 创建域名服务商数据结构体
type DomainServiceProviderCreate struct {
	Name      string  `json:"name" binding:"required"`
	AccessKey *string `json:"access_key"`
	SecretKey *string `json:"secret_key"`
	Type      uint    `json:"type" binding:"required"`
}

// AddDomainServiceProvider 创建域名服务商
func (d *domain) AddDomainServiceProvider(data *DomainServiceProviderCreate) (res *model.DomainServiceProvider, err error) {

	provider := &model.DomainServiceProvider{
		Name:      data.Name,
		Type:      data.Type,
		AccessKey: data.AccessKey,
		SecretKey: data.SecretKey,
	}

	return dao.Domain.AddDomainServiceProvider(provider)
}

// DeleteDomainServiceProvider 删除域名服务商
func (d *domain) DeleteDomainServiceProvider(id int) (err error) {
	return dao.Domain.DeleteDomainServiceProvider(id)
}

// UpdateDomainServiceProviderList 更新域名服务商
func (d *domain) UpdateDomainServiceProviderList(data *dao.ProviderUpdate) (*model.DomainServiceProvider, error) {

	updates := make(map[string]interface{})
	updates["id"] = data.ID
	updates["name"] = data.Name
	updates["type"] = data.Type

	if data.AccessKey != nil {
		// 加密
		ak, err := utils.Encrypt(*data.AccessKey)
		if err != nil {
			return nil, err
		}
		updates["access_key"] = ak
	} else {
		// 重置为空
		updates["access_key"] = nil
	}
	if data.SecretKey != nil {
		// 加密
		sk, err := utils.Encrypt(*data.SecretKey)
		if err != nil {
			return nil, err
		}
		updates["secret_key"] = sk
	} else {
		// 重置为空
		updates["secret_key"] = nil
	}

	return dao.Domain.UpdateDomainServiceProvider(updates)
}

// GetDomainServiceProviderList 获取域名服务商列表
func (d *domain) GetDomainServiceProviderList() ([]model.DomainServiceProvider, error) {
	return dao.Domain.GetDomainServiceProviderList()
}
