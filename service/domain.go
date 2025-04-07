package service

import (
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	"ops-api/config"
	"ops-api/dao"
	"ops-api/model"
	"ops-api/utils"
	"ops-api/utils/mail"
	"ops-api/utils/public_cloud"
	"strings"
	"time"
)

var Domain domain

type domain struct{}

// CloudProvider 云服务商相关接口
type CloudProvider interface {
	SyncDomains(serviceProviderID uint) ([]public_cloud.DomainList, error)
	GetDns(pageNum, pageSize int64, domainName, keyWord string) (*public_cloud.DnsList, error)
	AddDns(domainName, rrType, rr, value, remark string, ttl int32, weight *int32, priority int32) (recordId string, err error)
	UpdateDns(domainName, recordId, rrType, rr, value, remark string, ttl int32, weight *int32, priority int32) error
	DeleteDns(domainName, recordId string) error
	SetDnsStatus(domainName, recordId, status string) error
}

// DomainCreate 创建域名数据结构体
type DomainCreate struct {
	Name                    string     `json:"name" binding:"required"`
	RegistrationAt          *time.Time `json:"registration_at" binding:"required"`
	ExpirationAt            *time.Time `json:"expiration_at" binding:"required"`
	DomainServiceProviderID uint       `json:"domain_service_provider_id" binding:"required"`
}

// DomainServiceProviderCreate 创建域名服务商数据结构体
type DomainServiceProviderCreate struct {
	Name        string  `json:"name" binding:"required"`
	AccessKey   *string `json:"access_key"`
	SecretKey   *string `json:"secret_key"`
	Type        uint    `json:"type" binding:"required"`
	AutoSync    bool    `json:"auto_sync"`
	AccountName *string `json:"account_name"`
	IamUsername *string `json:"iam_username"`
	IamPassword *string `json:"iam_password"`
}

// DnsCreate 创建DNS记录结构体
type DnsCreate struct {
	DomainId uint   `json:"domain_id" binding:"required"`
	RR       string `json:"rr" binding:"required"`
	Type     string `json:"type" binding:"required"`
	TTL      int32  `json:"ttl" binding:"required"`
	Value    string `json:"value" binding:"required"`
	Priority int32  `json:"priority"`
	Weight   *int32 `json:"weight"`
	Remark   string `json:"remark"`
}

// DnsUpdate 修改DNS记录结构体
type DnsUpdate struct {
	DomainId uint   `json:"domain_id" binding:"required"`
	RecordId string `json:"record_id" binding:"required"`
	RR       string `json:"rr" binding:"required"`
	Type     string `json:"type" binding:"required"`
	Value    string `json:"value" binding:"required"`
	TTL      int32  `json:"ttl" binding:"required"`
	Priority int32  `json:"priority"`
	Weight   *int32 `json:"weight"`
	Remark   string `json:"remark"`
}

// DnsDelete 删除DNS记录结构体
type DnsDelete struct {
	DomainId uint   `json:"domain_id" binding:"required"`
	RecordId string `json:"record_id" binding:"required"`
}

// SetDnsStatus 设置DNS状态记录结构体
type SetDnsStatus struct {
	DomainId uint   `json:"domain_id" binding:"required"`
	RecordId string `json:"record_id" binding:"required"`
	Status   string `json:"status" binding:"required"`
}

// GetCloudProviderClient 获取云服务商客户端
func (d *domain) GetCloudProviderClient(provider *model.DomainServiceProvider) (CloudProvider, error) {

	if provider.AccessKey == nil || provider.SecretKey == nil {
		return nil, errors.New("服务商配置信息错误")
	}

	// 解密AccessKey和SecretKey
	ak, sk := decryptKeys(provider.AccessKey, provider.SecretKey)

	// 根据服务商类型创建客户端
	switch provider.Type {
	case 1:
		return public_cloud.CreateAliyunClient(ak, sk)
	case 2:
		return public_cloud.CreateHuaweiClient(ak, sk)
	case 3:
		return public_cloud.CreateTencentClient(ak, sk)
	default:
		return nil, errors.New("不支持的服务商类型")
	}
}

// AddDomainServiceProvider 创建域名服务商
func (d *domain) AddDomainServiceProvider(data *DomainServiceProviderCreate) (res *model.DomainServiceProvider, err error) {

	provider := &model.DomainServiceProvider{
		Name:      data.Name,
		Type:      data.Type,
		AccessKey: data.AccessKey,
		SecretKey: data.SecretKey,
		AutoSync:  data.AutoSync,
	}

	if data.Type == 2 {
		provider.AccountName = data.AccountName
		provider.IamUsername = data.IamUsername
		provider.IamPassword = data.IamPassword
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
	updates["auto_sync"] = data.AutoSync
	updates["account_name"] = data.AccountName
	updates["iam_username"] = data.IamUsername

	if data.AccessKey != nil {
		// 加密
		ak, err := utils.Encrypt(*data.AccessKey)
		if err != nil {
			return nil, err
		}
		updates["access_key"] = ak
	}
	if data.SecretKey != nil {
		// 加密
		sk, err := utils.Encrypt(*data.SecretKey)
		if err != nil {
			return nil, err
		}
		updates["secret_key"] = sk
	}
	if data.IamPassword != nil {
		// 加密
		ip, err := utils.Encrypt(*data.IamPassword)
		if err != nil {
			return nil, err
		}
		updates["iam_password"] = ip
	}

	return dao.Domain.UpdateDomainServiceProvider(updates)
}

// GetDomainServiceProviderList 获取域名服务商列表
func (d *domain) GetDomainServiceProviderList() ([]model.DomainServiceProvider, error) {
	return dao.Domain.GetDomainServiceProviderList()
}

// AddDomain 创建域名
func (d *domain) AddDomain(data *DomainCreate) (res *model.Domain, err error) {

	domain := &model.Domain{
		Name:                    data.Name,
		RegistrationAt:          data.RegistrationAt,
		ExpirationAt:            data.ExpirationAt,
		DomainServiceProviderID: data.DomainServiceProviderID,
	}

	return dao.Domain.AddDomain(domain)
}

// DeleteDomain 删除域名
func (d *domain) DeleteDomain(id int) (err error) {
	return dao.Domain.DeleteDomain(id)
}

// UpdateDomain 更新域名
func (d *domain) UpdateDomain(data *dao.DomainUpdate) (*model.Domain, error) {
	return dao.Domain.UpdateDomain(data)
}

// GetDomainList 获取域名列表
func (d *domain) GetDomainList(name string, providerId uint, page, limit *int) (data *dao.DomainList, err error) {
	return dao.Domain.GetDomainList(name, providerId, page, limit)
}

// SyncDomain 同步域名
func (d *domain) SyncDomain(ProviderId uint) error {

	// 获取域名服务商配置信息
	provider, err := dao.Domain.GetDomainServiceProviderForID(int(ProviderId))
	if err != nil {
		return err
	}

	// 创建请求客户端
	client, err := d.GetCloudProviderClient(provider)
	if err != nil {
		return err
	}

	// 获取域名列表
	domains, err := client.SyncDomains(provider.Id)
	if err != nil {
		return err
	}

	// 将获取到的域名列表转换为 []*model.Domain
	var modelDomains []*model.Domain
	for _, d := range domains {
		modelDomains = append(modelDomains, &model.Domain{
			Name:                    d.Name,
			RegistrationAt:          d.RegistrationAt,
			ExpirationAt:            d.ExpirationAt,
			DomainServiceProviderID: d.DomainServiceProviderID,
		})
	}

	return dao.Domain.SyncDomains(modelDomains, provider.Id)
}

// GetDnsList 获取域名DNS解析列表
func (d *domain) GetDnsList(keyWord string, ID uint, page, limit int) (*public_cloud.DnsList, error) {
	// 获取域名信息
	result, err := dao.Domain.GetDomainForID(ID)
	if err != nil {
		return nil, err
	}

	// 获取域名服务商配置信息
	provider, err := dao.Domain.GetDomainServiceProviderForID(int(result.DomainServiceProviderID))
	if err != nil {
		return nil, err
	}

	// 创建请求客户端
	client, err := d.GetCloudProviderClient(provider)
	if err != nil {
		return nil, err
	}

	data, err := client.GetDns(int64(page), int64(limit), result.Name, keyWord)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// AddDns 新增域名DNS解析
func (d *domain) AddDns(dns *DnsCreate) error {
	// 获取域名信息
	result, err := dao.Domain.GetDomainForID(dns.DomainId)
	if err != nil {
		return err
	}

	// 获取域名服务商配置信息
	provider, err := dao.Domain.GetDomainServiceProviderForID(int(result.DomainServiceProviderID))
	if err != nil {
		return err
	}

	// 创建请求客户端
	client, err := d.GetCloudProviderClient(provider)
	if err != nil {
		return err
	}

	_, err = client.AddDns(result.Name, dns.Type, dns.RR, dns.Value, dns.Remark, dns.TTL, dns.Weight, dns.Priority)
	return err
}

// UpdateDomainDns 修改域名DNS解析
func (d *domain) UpdateDomainDns(dns *DnsUpdate) error {
	// 获取域名信息
	result, err := dao.Domain.GetDomainForID(dns.DomainId)
	if err != nil {
		return err
	}

	// 获取域名服务商配置信息
	provider, err := dao.Domain.GetDomainServiceProviderForID(int(result.DomainServiceProviderID))
	if err != nil {
		return err
	}

	// 创建请求客户端
	client, err := d.GetCloudProviderClient(provider)
	if err != nil {
		return err
	}

	return client.UpdateDns(result.Name, dns.RecordId, dns.Type, dns.RR, dns.Value, dns.Remark, dns.TTL, dns.Weight, dns.Priority)
}

// DeleteDns 删除域名DNS解析
func (d *domain) DeleteDns(dns *DnsDelete) error {
	// 获取域名信息
	result, err := dao.Domain.GetDomainForID(dns.DomainId)
	if err != nil {
		return err
	}

	// 获取域名服务商配置信息
	provider, err := dao.Domain.GetDomainServiceProviderForID(int(result.DomainServiceProviderID))
	if err != nil {
		return err
	}

	// 创建请求客户端
	client, err := d.GetCloudProviderClient(provider)
	if err != nil {
		return err
	}

	return client.DeleteDns(result.Name, dns.RecordId)
}

// SetDnsStatus 设置域名DNS状态
func (d *domain) SetDnsStatus(dns *SetDnsStatus) error {
	// 获取域名信息
	result, err := dao.Domain.GetDomainForID(dns.DomainId)
	if err != nil {
		return err
	}

	// 获取域名服务商配置信息
	provider, err := dao.Domain.GetDomainServiceProviderForID(int(result.DomainServiceProviderID))
	if err != nil {
		return err
	}

	// 创建请求客户端
	client, err := d.GetCloudProviderClient(provider)
	if err != nil {
		return err
	}

	return client.SetDnsStatus(result.Name, dns.RecordId, dns.Status)
}

// DomainExpiredNotice 域名过期通知
func (d *domain) DomainExpiredNotice() error {

	domains, err := dao.Domain.GetExpiredDomainList()
	if err != nil {
		return err
	}

	if len(domains) == 0 {
		logger.Info("检查域名状态正常.")
		return nil
	}

	// 生成HTML内容
	htmlBody := domainExpiredNoticeHTML(domains)

	// 发送邮件函数
	return mail.Email.SendMsg([]string{"270142877@qq.com"}, nil, nil, "域名过期提醒", htmlBody, "html")
}

// domainExpiredNoticeHTML 域名过期通知 HTML
func domainExpiredNoticeHTML(domains []*model.Domain) string {

	var (
		issuer = config.Conf.Settings["issuer"].(string)
		now    = time.Now()
	)

	return fmt.Sprintf(`
        <!DOCTYPE html>
        <html>
        <head>
            <meta charset="UTF-8">
            <title>域名过期提醒</title>
            <style>
                /* 主容器设置固定宽度并居中 */
                .email-container {
                    max-width: 800px;
                    margin: 0 auto;
                    font-family: Arial, sans-serif;
                    line-height: 1.6;
                    color: #333;
                    padding: 20px;
                }
                
                /* 内容区域样式 */
                .email-content {
                    background-color: #ffffff;
                    border: 1px solid #e0e0e0;
                    border-radius: 4px;
                    padding: 25px;
                }
                
                h1 {
                    color: #2c3e50;
                    border-bottom: 2px solid #3498db;
                    padding-bottom: 10px;
                    margin-top: 0;
                }
                
                .info {
                    background-color: #f8f9fa;
                    padding: 15px;
                    border-left: 4px solid #3498db;
                    margin-bottom: 20px;
                }
                
                table {
                    width: 100%%;
                    border-collapse: collapse;
                    margin: 20px 0;
                }
                
                th {
                    background-color: #3498db;
                    color: white;
                    padding: 12px;
                    text-align: left;
                }
                
                td {
                    padding: 10px;
                    border-bottom: 1px solid #ddd;
                }
                
                tr:nth-child(even) {
                    background-color: #f2f2f2;
                }
                
                .expired {
                    color: #e74c3c;
                    font-weight: bold;
                }
                
                .expiring {
                    color: #f39c12;
                    font-weight: bold;
                }
                
                .footer {
                    margin-top: 30px;
                    font-size: 0.9em;
                    color: #7f8c8d;
                    /* 移除 text-align: center 改为默认左对齐 */
                }
                
                .auto-send-notice {
                    color: #e74c3c;
                    margin-top: 10px;
                }
            </style>
        </head>
        <body>
            <!-- 主容器 -->
            <div class="email-container">
                <!-- 内容区域 -->
                <div class="email-content">
                    <h1>域名过期提醒</h1>
                    <div class="info">
                        以下域名即将过期或已过期，请及时处理以避免服务中断。
                    </div>
                    
                    <table>
                        <thead>
                            <tr>
                                <th>域名</th>
								<th>域名服务商</th>
                                <th>过期时间</th>
                                <th>状态</th>
                            </tr>
                        </thead>
                        <tbody>
                            %s
                        </tbody>
                    </table>
                    
                    <div class="footer">
                        <p>此致,<br>%s</p>
                        <p class="auto-send-notice">此邮件为系统自动发送，请不要回复此邮件。</p>
                    </div>
                </div>
            </div>
        </body>
        </html>
    `, generateDomainRows(domains, now), issuer)
}

// generateDomainRows 生成域名表格行
func generateDomainRows(domains []*model.Domain, now time.Time) string {
	var rows strings.Builder

	for _, domain := range domains {

		// 确定域名状态
		status := "正常"
		statusClass := ""

		if domain.ExpirationAt.Before(now) {
			status = "已过期"
			statusClass = "expired"
		} else if domain.ExpirationAt.Before(now.Add(30 * 24 * time.Hour)) {
			status = "即将过期"
			statusClass = "expiring"
		}

		// 格式化过期时间
		expiredAt := domain.ExpirationAt.Format("2006-01-02 15:04:05")

		rows.WriteString(fmt.Sprintf(`
            <tr>
                <td>%s</td>
				<td>%s</td>
                <td>%s</td>
                <td class="%s">%s</td>
            </tr>
        `, domain.Name, domain.DomainServiceProvider.Name, expiredAt, statusClass, status))
	}

	return rows.String()
}

// decryptKeys 解密 AccessKey 和 SecretKey
func decryptKeys(accessKey, secretKey *string) (string, string) {
	ak, _ := utils.Decrypt(*accessKey)
	sk, _ := utils.Decrypt(*secretKey)
	return ak, sk
}
