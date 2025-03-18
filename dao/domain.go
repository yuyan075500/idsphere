package dao

import (
	"errors"
	"gorm.io/gorm"
	"ops-api/global"
	"ops-api/model"
	"ops-api/utils"
	"time"
)

var Domain domain

type domain struct{}

// DomainList 返回给前端表格的数据结构体
type DomainList struct {
	Items []*model.Domain `json:"items"`
	Total int64           `json:"total"`
}

// DomainUpdate 更新域名结构体
type DomainUpdate struct {
	ID                      uint       `json:"id" binding:"required"`
	Name                    string     `json:"name" binding:"required"`
	RegistrationAt          *time.Time `json:"registration_at" binding:"required"`
	ExpirationAt            *time.Time `json:"expiration_at" binding:"required"`
	DomainServiceProviderID uint       `json:"domain_service_provider_id" binding:"required"`
}

// ProviderUpdate 更新域名服务商结构体
type ProviderUpdate struct {
	ID          uint    `json:"id" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	AccessKey   *string `json:"access_key"`
	SecretKey   *string `json:"secret_key"`
	Type        uint    `json:"type" binding:"required"`
	AutoSync    bool    `json:"auto_sync"`
	AccountName *string `json:"account_name"`
	IamUsername *string `json:"iam_username"`
	IamPassword *string `json:"iam_password"`
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
	if err := global.MySQLClient.Where("id = ?", id).Unscoped().Delete(&model.DomainServiceProvider{}).Error; err != nil {
		if utils.IsForeignKeyConstraintError(err) {
			return errors.New("请确保服务商下不存在域名")
		}
		return err
	}
	return nil
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
	if err := global.MySQLClient.Omit("AccessKey", "SecretKey", "IamPassword").Find(&providers).Error; err != nil {
		return nil, err
	}

	return providers, nil
}

// GetDomainServiceProviderForID 根据 ID 获取域名服务商
func (d *domain) GetDomainServiceProviderForID(id int) (*model.DomainServiceProvider, error) {
	var provider model.DomainServiceProvider
	if err := global.MySQLClient.First(&provider, id).Error; err != nil {
		return nil, err
	}
	return &provider, nil
}

// AddDomain 新增域名
func (d *domain) AddDomain(data *model.Domain) (provider *model.Domain, err error) {
	if err := global.MySQLClient.Create(&data).Error; err != nil {
		return nil, err
	}
	return provider, nil
}

// SyncDomains 同步域名
func (d *domain) SyncDomains(domains []*model.Domain, providerId uint) error {
	return global.MySQLClient.Transaction(func(tx *gorm.DB) error {

		for _, result := range domains {
			// 查找当前服务商下是否已存在相同的域名
			var existing model.Domain
			if err := tx.Where("name = ? AND domain_service_provider_id = ?", result.Name, result.DomainServiceProviderID).First(&existing).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// 不存在则创建
					if err := tx.Create(result).Error; err != nil {
						return err
					}
				} else {
					return err
				}
			} else {
				// 存在则更新
				if err := tx.Model(&existing).Updates(map[string]interface{}{
					"registration_at": result.RegistrationAt,
					"expiration_at":   result.ExpirationAt,
				}).Error; err != nil {
					return err
				}
			}
		}

		// 生成一个所有域名的名称列表
		var domainNames []string
		for _, result := range domains {
			domainNames = append(domainNames, result.Name)
		}

		// 删除数据库中多余的域名记录
		if len(domainNames) > 0 {
			if err := tx.Where("name NOT IN (?) AND domain_service_provider_id = ?", domainNames, providerId).Delete(&model.Domain{}).Error; err != nil {
				return err
			}
		} else {
			// 如果为空，仅删除当前服务商的数据
			if err := tx.Where("domain_service_provider_id = ?", providerId).Delete(&model.Domain{}).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// DeleteDomain 删除域名
func (d *domain) DeleteDomain(id int) (err error) {
	return global.MySQLClient.Where("id = ?", id).Unscoped().Delete(&model.Domain{}).Error
}

// UpdateDomain 修改域名
func (d *domain) UpdateDomain(data *DomainUpdate) (*model.Domain, error) {

	domain := &model.Domain{}

	if err := global.MySQLClient.Model(domain).Select("*").Where("id = ?", data.ID).Updates(data).Error; err != nil {
		return nil, err
	}

	// 查询更新后的账号信息并返回
	if err := global.MySQLClient.First(domain, data.ID).Error; err != nil {
		return nil, err
	}
	return domain, nil
}

// GetDomainList 获取域名列表
func (d *domain) GetDomainList(name string, providerId uint, page, limit int) (*DomainList, error) {
	var (
		startSet = (page - 1) * limit
		domains  []*model.Domain
		total    int64
	)

	// 初始化查询
	query := global.MySQLClient.Model(&model.Domain{}).
		Preload("DomainServiceProvider", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).Where("name LIKE ?", "%"+name+"%")

	// 过滤服务商
	if providerId != 0 {
		query = query.Where("domain_service_provider_id = ?", providerId)
	}

	// 执行查询
	if err := query.Count(&total).Limit(limit).Offset(startSet).Find(&domains).Error; err != nil {
		return nil, err
	}

	return &DomainList{
		Items: domains,
		Total: total,
	}, nil
}

// GetDomainForID 根据 ID 获取域名
func (d *domain) GetDomainForID(id uint) (*model.Domain, error) {
	var result model.Domain
	if err := global.MySQLClient.Preload("DomainServiceProvider").First(&result, id).Error; err != nil {
		return nil, err
	}
	return &result, nil
}
