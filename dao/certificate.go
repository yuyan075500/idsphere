package dao

import (
	"ops-api/global"
	"ops-api/model"
)

var Certificate certificate

type certificate struct{}

// DomainCertificateList 返回给前端表格的证书列表
type DomainCertificateList struct {
	Items []*model.DomainCertificate `json:"items"`
	Total int64                      `json:"total"`
}

// DomainCertificateRequestList 返回给前端表格的证书申请列表
type DomainCertificateRequestList struct {
	Items []*model.DomainCertificateRequestRecord `json:"items"`
	Total int64                                   `json:"total"`
}

// CreateDomainCertificateRequest 创建证书申请记录
func (c *certificate) CreateDomainCertificateRequest(data *model.DomainCertificateRequestRecord) (request *model.DomainCertificateRequestRecord, err error) {
	if err := global.MySQLClient.Create(&data).Error; err != nil {
		return nil, err
	}
	return request, nil
}

// GetCertificateRequestForID 根据 ID 获取证书申请记录
func (c *certificate) GetCertificateRequestForID(id int) (*model.DomainCertificateRequestRecord, error) {
	var record model.DomainCertificateRequestRecord
	if err := global.MySQLClient.First(&record, id).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

// UploadDomainCertificate 上传证书
func (c *certificate) UploadDomainCertificate(data *model.DomainCertificate) (provider *model.DomainCertificate, err error) {
	if err := global.MySQLClient.Create(&data).Error; err != nil {
		return nil, err
	}
	return provider, nil
}

// DeleteCertificate 删除证书
func (c *certificate) DeleteCertificate(id int) (err error) {
	if err := global.MySQLClient.Where("id = ?", id).Unscoped().Delete(&model.DomainCertificate{}).Error; err != nil {
		return err
	}
	return nil
}

// GetDomainCertificateList 获取证书列表
func (c *certificate) GetDomainCertificateList(name string, page, limit int) (*DomainCertificateList, error) {
	var (
		startSet = (page - 1) * limit
		certs    []*model.DomainCertificate
		total    int64
	)

	// 初始化查询
	query := global.MySQLClient.Model(&model.DomainCertificate{}).Where("domain LIKE ?", "%"+name+"%")

	// 执行查询
	if err := query.Omit("certificate", "private_key").
		Count(&total).Limit(limit).
		Offset(startSet).
		Find(&certs).Error; err != nil {
		return nil, err
	}

	return &DomainCertificateList{
		Items: certs,
		Total: total,
	}, nil
}

// GetCertificateForID 根据 ID 获取证书
func (c *certificate) GetCertificateForID(id int) (*model.DomainCertificate, error) {
	var cert model.DomainCertificate
	if err := global.MySQLClient.First(&cert, id).Error; err != nil {
		return nil, err
	}
	return &cert, nil
}
