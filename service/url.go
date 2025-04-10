package service

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	"net"
	"ops-api/config"
	"ops-api/dao"
	"ops-api/global"
	"ops-api/model"
	"ops-api/utils/mail"
	"strconv"
	"strings"
	"time"
)

var UrlAddress urlAddress

type urlAddress struct{}

// UrlAddressCreate 创建数据结构体
type UrlAddressCreate struct {
	Name      string `json:"name" binding:"required"`
	Domain    string `json:"domain" binding:"required"`
	Port      uint   `json:"port" binding:"required"`
	IPAddress string `json:"ip_address"`
}

// AddUrl 创建
func (u *urlAddress) AddUrl(data *UrlAddressCreate) (res *model.DomainCertificateMonitor, err error) {

	url := &model.DomainCertificateMonitor{
		Domain:    data.Domain,
		Name:      data.Name,
		Port:      data.Port,
		IPAddress: data.IPAddress,
	}

	return dao.UrlAddress.AddUrl(url)
}

// DeleteUrl 删除
func (u *urlAddress) DeleteUrl(id int) (err error) {
	return dao.UrlAddress.DeleteUrl(id)
}

// UpdateUrl 更新
func (u *urlAddress) UpdateUrl(data *dao.UrlAddressUpdate) (*model.DomainCertificateMonitor, error) {
	return dao.UrlAddress.UpdateUrl(data)
}

// GetUrlList 获取列表
func (u *urlAddress) GetUrlList(name string, page, limit *int) (data *dao.UrlAddressList, err error) {
	return dao.UrlAddress.GetUrlList(name, page, limit)
}

// CertificateCheck 证书检查
func (u *urlAddress) CertificateCheck(id *uint) error {

	var ids []uint

	if id == nil {
		// 获取所有 UrlAddress
		urls, err := dao.UrlAddress.GetUrlList("", nil, nil)
		if err != nil {
			return err
		}
		for _, url := range urls.Items {
			ids = append(ids, url.ID)
		}

		// 检查证书
		for _, checkID := range ids {
			if err := u.checkSingleCertificate(checkID); err != nil {
				continue
			}
		}

		// 获取所有状态总异常的记录
		records, err := dao.UrlAddress.GetExpiredOrExpirationList()
		if err != nil {
			return err
		}

		// 生成HTML内容
		htmlBody := urlCertificateExpiredNoticeHTML(records)

		// 发送邮件告警
		return mail.Email.SendMsg([]string{"270142877@qq.com"}, nil, nil, "URL 站点证书过期提醒", htmlBody, "html")
	} else {
		// 检查证书
		if err := u.checkSingleCertificate(*id); err != nil {
			logger.Error(err)
		}
	}

	return nil
}

// urlCertificateExpiredNoticeHTML URL证书过期通知 HTML
func urlCertificateExpiredNoticeHTML(urls []*model.DomainCertificateMonitor) string {

	var (
		issuer = config.Conf.Settings["issuer"].(string)
		now    = time.Now()
	)

	return fmt.Sprintf(`
        <!DOCTYPE html>
        <html>
        <head>
            <meta charset="UTF-8">
            <title>URL 站点证书过期提醒</title>
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

				.text-muted {
					color: #888;
					font-style: italic;
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
                    <h1>URL 站点证书过期提醒</h1>
                    <div class="info">
                        以下URL 站点证书即将过期或已过期，请及时处理以避免服务中断。
                    </div>
                    
                    <table>
                        <thead>
                            <tr>
                                <th>站点名称</th>
								<th>域名</th>
                                <th>过期时间</th>
								<th>检查时间</th>
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
    `, generateUrlCertificateRows(urls, now), issuer)
}

// generateUrlCertificateRows 生成表格行
func generateUrlCertificateRows(urls []*model.DomainCertificateMonitor, now time.Time) string {
	var rows strings.Builder

	for _, url := range urls {
		var (
			statusText  = "未检查"
			statusClass = "text-muted"
			expiredAt   = "-"
			checkedAt   = "-"
		)

		// 判断状态
		if url.Status != nil {
			switch *url.Status {
			case 0:
				statusText = "正常"
				statusClass = ""
				if url.ExpirationAt != nil && url.ExpirationAt.Before(now.Add(90*24*time.Hour)) {
					statusText = "即将过期"
					statusClass = "expiring"
				}
			case 1:
				statusText = "检查异常"
				statusClass = "expired"
			case 2:
				statusText = "已过期"
				statusClass = "expired"
			}
		}

		// 格式化时间
		if url.ExpirationAt != nil {
			expiredAt = url.ExpirationAt.Format("2006-01-02 15:04:05")
		}
		if url.LastCheckAt != nil {
			checkedAt = url.LastCheckAt.Format("2006-01-02 15:04:05")
		}

		rows.WriteString(fmt.Sprintf(`
            <tr>
                <td>%s</td>
                <td>%s</td>
                <td>%s</td>
                <td>%s</td>
				<td class="%s">%s</td>
            </tr>
        `, url.Name, url.Domain, expiredAt, checkedAt, statusClass, statusText))
	}

	return rows.String()
}

// 检查单个证书
func (u *urlAddress) checkSingleCertificate(id uint) error {
	status := 0

	url, err := dao.UrlAddress.GetUrlForID(id)
	if err != nil {
		return fmt.Errorf("数据库读取信息失败: %w", err)
	}

	// 获取证书
	cert, err := fetchCertificate(url.Domain, url.IPAddress, url.Port)
	if err != nil {
		logger.Error("证书获取失败 (Domain: %v): %v", url.Domain, err.Error())
		status = 1
		if err := updateStatus(&id, status); err != nil {
			return fmt.Errorf("状态更新失败: %w", err)
		}
		return nil
	}

	if cert.NotAfter.Before(time.Now()) {
		status = 2
	}

	if err := updateStatusAndExpiration(&id, status, cert.NotAfter); err != nil {
		return fmt.Errorf("证书信息更新失败: %w", err)
	}

	return nil
}

// updateStatusAndExpiration 更新证书状态和过期时间
func updateStatusAndExpiration(id *uint, status int, expirationTime time.Time) error {
	now := time.Now()
	urlAddress := &model.DomainCertificateMonitor{
		ID:           *id,
		Status:       &status,
		ExpirationAt: &expirationTime,
		LastCheckAt:  &now,
	}
	return global.MySQLClient.Model(&model.DomainCertificateMonitor{}).Where("id = ?", *id).Select("Status", "ExpirationAt", "LastCheckAt").Updates(urlAddress).Error
}

// updateStatus 更新证书状态
func updateStatus(id *uint, status int) error {
	urlAddress := &model.DomainCertificateMonitor{
		ID:     *id,
		Status: &status,
	}
	return global.MySQLClient.Model(&model.DomainCertificateMonitor{}).
		Where("id = ?", *id).
		Select("Status").
		Updates(urlAddress).Error
}

func fetchCertificate(domain, ip string, port uint) (*x509.Certificate, error) {

	var (
		address    = net.JoinHostPort(domain, strconv.Itoa(int(port)))
		serverName = domain
		dialer     = &net.Dialer{
			Timeout: 5 * time.Second,
		}
	)

	if ip != "" {
		address = net.JoinHostPort(ip, strconv.Itoa(int(port)))
	}

	conn, err := tls.DialWithDialer(dialer, "tcp", address, &tls.Config{
		ServerName:         serverName,
		InsecureSkipVerify: true,
	})
	if err != nil {
		return nil, err
	}

	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return nil, errors.New("no certificates found")
	}

	return certs[0], nil
}
