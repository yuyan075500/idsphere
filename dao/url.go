package dao

import (
	"ops-api/global"
	"ops-api/model"
	"time"
)

var UrlAddress urlAddress

type urlAddress struct{}

// UrlAddressList 返回给前端表格的数据结构体
type UrlAddressList struct {
	Items []*model.DomainCertificateMonitor `json:"items"`
	Total int64                             `json:"total"`
}

// UrlAddressUpdate 更新结构体
type UrlAddressUpdate struct {
	ID        uint   `json:"id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Domain    string `json:"domain" binding:"required"`
	Port      uint   `json:"port" binding:"required"`
	IPAddress string `json:"ip_address"`
}

// AddUrl 新增
func (u *urlAddress) AddUrl(data *model.DomainCertificateMonitor) (url *model.DomainCertificateMonitor, err error) {
	if err := global.MySQLClient.Create(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// DeleteUrl 删除
func (u *urlAddress) DeleteUrl(id int) (err error) {
	if err := global.MySQLClient.Where("id = ?", id).Unscoped().Delete(&model.DomainCertificateMonitor{}).Error; err != nil {
		return err
	}
	return nil
}

// UpdateUrl 修改
func (u *urlAddress) UpdateUrl(data *UrlAddressUpdate) (*model.DomainCertificateMonitor, error) {

	url := &model.DomainCertificateMonitor{}

	if err := global.MySQLClient.Model(url).Select("*").Where("id = ?", data.ID).Updates(data).Error; err != nil {
		return nil, err
	}

	// 查询更新后的信息并返回
	if err := global.MySQLClient.First(url, data.ID).Error; err != nil {
		return nil, err
	}
	return url, nil
}

// GetUrlForID 根据 ID 获取
func (u *urlAddress) GetUrlForID(id uint) (*model.DomainCertificateMonitor, error) {
	var url model.DomainCertificateMonitor
	if err := global.MySQLClient.First(&url, id).Error; err != nil {
		return nil, err
	}
	return &url, nil
}

// GetUrlList 获取列表
func (u *urlAddress) GetUrlList(name string, page, limit *int) (*UrlAddressList, error) {
	var (
		url   []*model.DomainCertificateMonitor
		total int64
	)

	// 初始化查询
	query := global.MySQLClient.Model(&model.DomainCertificateMonitor{}).Where("name LIKE ? OR domain LIKE ?", "%"+name+"%", "%"+name+"%")

	if page == nil && limit == nil {
		// 不分页
		if err := query.Count(&total).Find(&url).Error; err != nil {
			return nil, err
		}
	} else {
		// 分页
		startSet := (*page - 1) * *limit
		if err := query.Count(&total).Limit(*limit).Offset(startSet).Find(&url).Error; err != nil {
			return nil, err
		}
	}

	return &UrlAddressList{
		Items: url,
		Total: total,
	}, nil
}

// GetExpiredOrExpirationList 获取证书过期或异常的站点列表
func (u *urlAddress) GetExpiredOrExpirationList() (urlList []*model.DomainCertificateMonitor, err error) {

	var (
		urls           []*model.DomainCertificateMonitor
		now            = time.Now()
		sevenDaysLater = now.Add(30 * 24 * time.Hour)
	)

	if err := global.MySQLClient.Model(&model.DomainCertificateMonitor{}).
		Where("status != ? OR expiration_at < ? OR expiration_at BETWEEN ? AND ?", 0, now, now, sevenDaysLater).
		Find(&urls).
		Error; err != nil {
		return nil, err
	}

	return urls, nil
}
