package service

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	"net"
	"ops-api/config"
	"ops-api/dao"
	"ops-api/global"
	"ops-api/model"
	"ops-api/utils/notify"
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

// ManualCertificateCheck 手动检查
func (u *urlAddress) ManualCertificateCheck(id *uint) error {
	url, err := dao.UrlAddress.GetUrlForID(*id)
	if err != nil {
		return err
	}
	return u.checkSingleCertificate(url)
}

// AutoCertificateCheck 定量任务检查
func (u *urlAddress) AutoCertificateCheck(task *model.ScheduledTask) error {

	// 获取所有 URL 站点
	urls, err := dao.UrlAddress.GetUrlList("", nil, nil)
	if err != nil {
		return err
	}

	// 循环执行检查
	for _, url := range urls.Items {
		if err := u.checkSingleCertificate(url); err != nil {
			logger.Error("证书验证失败，域名：%s，ErrorMsg：%s", url.Domain, err.Error())
			continue
		}
	}

	// 获取所有状态为异常地记录（过期、检查异常、即将过期）
	records, err := dao.UrlAddress.GetExpiredOrExpirationList()
	if err != nil {
		return err
	}

	// 如果记录为空则不做任何操作
	if len(records) < 0 {
		logger.Info("检查 URL 站点状态正常.")
		return nil
	}

	// 生成通知内容（1：邮件 HTML，3：富文本，其它： Markdown 文档）
	notifyType := *task.NotifyType
	var message string
	switch notifyType {
	case 1:
		message = urlCertificateExpiredNoticeHTML(records)
	case 3:
		postData := urlCertificateExpiredNoticeFeishuPost(records)
		jsonBytes, _ := json.Marshal(postData)
		message = string(jsonBytes)
	default:
		message = urlCertificateExpiredNoticeMarkdown(records)
	}

	// 发送告警
	notifier := notify.GetNotifier(*task)
	return notifier.SendNotify(message, "URL 站点 HTTPS 证书过期提醒")
}

// urlCertificateExpiredNoticeFeishuPost 生成飞书 Post 富文本消息
func urlCertificateExpiredNoticeFeishuPost(urls []*model.DomainCertificateMonitor) map[string]interface{} {
	var (
		now     = time.Now()
		issuer  = config.Conf.Settings["issuer"].(string)
		content = make([][]map[string]interface{}, 0)
	)

	for i, url := range urls {
		var (
			statusText  = "未检查"
			statusColor = "default"
			expiredAt   = "-"
			checkedAt   = "-"
		)

		if url.ExpirationAt != nil {
			expiredAt = url.ExpirationAt.Format("2006-01-02 15:04:05")
		}
		if url.LastCheckAt != nil {
			checkedAt = url.LastCheckAt.Format("2006-01-02 15:04:05")
		}

		if url.Status != nil {
			switch *url.Status {
			case 0:
				if url.ExpirationAt != nil && url.ExpirationAt.Before(now.Add(90*24*time.Hour)) {
					statusText = "即将过期"
					statusColor = "orange"
				} else {
					statusText = "正常"
					statusColor = "green"
				}
			case 1:
				statusText = "检查异常"
				statusColor = "orange"
			case 2:
				statusText = "已过期"
				statusColor = "red"
			}
		}

		// 每条记录拼接一段
		content = append(content, []map[string]interface{}{
			{"tag": "text", "text": fmt.Sprintf("%d. 名称：", i+1)},
			{"tag": "text", "text": url.Name, "bold": true},
		})
		content = append(content, []map[string]interface{}{
			{"tag": "text", "text": "   域名："},
			{"tag": "a", "text": url.Domain},
		})
		content = append(content, []map[string]interface{}{
			{"tag": "text", "text": "   到期时间："},
			{"tag": "text", "text": expiredAt},
		})
		content = append(content, []map[string]interface{}{
			{"tag": "text", "text": "   检查时间："},
			{"tag": "text", "text": checkedAt},
		})
		content = append(content, []map[string]interface{}{
			{"tag": "text", "text": "   状态："},
			{"tag": "text", "text": statusText, "text_color": statusColor},
		})
	}

	content = append(content, []map[string]interface{}{
		{"tag": "text", "text": "--------------------------------\n"},
	})
	content = append(content, []map[string]interface{}{
		{"tag": "text", "text": fmt.Sprintf("来源：%s", issuer)},
	})

	return map[string]interface{}{
		"msg_type": "post",
		"content": map[string]interface{}{
			"post": map[string]interface{}{
				"zh_cn": map[string]interface{}{
					"title":   "URL 站点 HTTPS 证书异常提醒：",
					"content": content,
				},
			},
		},
	}
}

// urlCertificateExpiredNoticeMarkdown 生成 URL HTTPS 证书过期通知 Markdown 文档
func urlCertificateExpiredNoticeMarkdown(urls []*model.DomainCertificateMonitor) string {
	var (
		builder = &strings.Builder{}
		now     = time.Now()
		issuer  = config.Conf.Settings["issuer"].(string)
	)

	builder.WriteString("**URL 站点 HTTPS 证书异常提醒：**\n\n")

	for i, url := range urls {
		var (
			statusText = "未检查"
			expiredAt  = "-"
			checkedAt  = "-"
		)

		if url.ExpirationAt != nil {
			expiredAt = url.ExpirationAt.Format("2006-01-02 15:04:05")
		}
		if url.LastCheckAt != nil {
			checkedAt = url.LastCheckAt.Format("2006-01-02 15:04:05")
		}

		if url.Status != nil {
			switch *url.Status {
			case 0:
				if url.ExpirationAt != nil && url.ExpirationAt.Before(now.Add(30*24*time.Hour)) {
					statusText = "<font color=\"warning\">即将过期</font>"
				} else {
					statusText = "<font color=\"info\">检查异常</font>"
				}
			case 1:
				statusText = "<font color=\"warning\">检查异常</font>"
			case 2:
				statusText = "<font color=\"warning\">已过期</font>"
			}
		}

		builder.WriteString(fmt.Sprintf("%d. 名称：%s\n\n", i+1, url.Name))
		builder.WriteString(fmt.Sprintf("   域名：%s\n\n", url.Domain))
		builder.WriteString(fmt.Sprintf("   到期时间：%s\n\n", expiredAt))
		builder.WriteString(fmt.Sprintf("   检查时间：%s\n\n", checkedAt))
		builder.WriteString(fmt.Sprintf("   状态：%s\n\n", statusText))
	}

	builder.WriteString("--------------------------------\n")
	builder.WriteString(fmt.Sprintf("来源：%s\n", issuer))

	return builder.String()
}

// urlCertificateExpiredNoticeHTML 生成 URL HTTPS 证书过期通知 HTML 文档
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
            <title>URL 站点 HTTPS 证书异常提醒</title>
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
                    <h1>URL 站点 HTTPS 证书异常提醒</h1>
                    <div class="info">
                        以下 URL 站点 HTTPS 证书异常，请及时处理以避免服务中断。
                    </div>
                    
                    <table>
                        <thead>
                            <tr>
                                <th>名称</th>
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

// generateUrlCertificateRows URL HTTPS 证书过期通知 HTML 表格数据渲染
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

// checkSingleCertificate 检查证书
func (u *urlAddress) checkSingleCertificate(data *model.DomainCertificateMonitor) error {

	// 向 URL 发起请求，获取证书
	cert, err := fetchCertificate(data.Domain, data.IPAddress, data.Port)
	if err != nil {
		return err
	}

	// 设置最后检查时间
	data.LastCheckAt = ptr(time.Now())

	if err != nil {
		// 设置状态为：1（表示异常）
		data.Status = ptr(1)

		// 在数据库中更新状态
		return global.MySQLClient.Save(data).Error
	}

	// 设置证书过期时间
	data.ExpirationAt = &cert.NotAfter

	// 判断证书是否过期
	if cert.NotAfter.Before(time.Now()) {
		// 设置状态为：2（表示已过期）
		data.Status = ptr(2)
	} else {
		// 设置状态为：0（表示正常）
		data.Status = ptr(0)
	}

	return global.MySQLClient.Save(data).Error
}

// ptr 创建一个指针
func ptr[T any](v T) *T {
	return &v
}

// fetchCertificate 从指定的 URL 地址获取证书信息
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
