package service

import (
	"archive/zip"
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/go-acme/lego/v4/certcrypto"
	cert "github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"github.com/wonderivan/logger"
	"ops-api/config"
	"ops-api/dao"
	"ops-api/global"
	"ops-api/model"
	"ops-api/utils/mail"
	"os"
	"strings"
	"time"
)

var Certificate certificate

type certificate struct{}

// DomainCertificateCreate 创建证书结构体
type DomainCertificateCreate struct {
	Certificate  string    `json:"certificate" binding:"required"`
	PrivateKey   string    `json:"private_key" binding:"required"`
	Domain       string    `json:"domain"`
	Type         uint      `json:"type" binding:"required"`
	ServerType   uint      `json:"server_type" binding:"required"`
	StartAt      time.Time `json:"start_at"`
	ExpirationAt time.Time `json:"expiration_at"`
	Status       string    `json:"status"`
}

// DomainCertificateRequest 证书申请
type DomainCertificateRequest struct {
	Email      *string `json:"email"`
	Domain     string  `json:"domain" binding:"required"`
	RR         string  `json:"rr" binding:"required"`
	ProviderID uint    `json:"provider_id" binding:"required"`
}

type AcmeUser struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *AcmeUser) GetEmail() string {
	return u.Email
}
func (u AcmeUser) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *AcmeUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

type dnsProvider struct {
	providerClient CloudProvider
	domain         string
	recordIds      map[string]string
}

func (p *dnsProvider) Present(domain, token, keyAuth string) error {
	// 获取DNS记录信息
	info := dns01.GetChallengeInfo(domain, keyAuth)

	// 解析当前这个域名的子域名部分
	sub := strings.TrimSuffix(domain, "."+p.domain)

	// 添加DNS记录
	rr := fmt.Sprintf("_acme-challenge.%s", sub)
	if sub == domain {
		rr = "_acme-challenge"
	}
	recordId, err := p.providerClient.AddDns(p.domain, "TXT", rr, info.Value, "DNS-01挑战", 600, nil, 0)
	if err != nil {
		return err
	}

	// 初始化记录IDs
	if p.recordIds == nil {
		p.recordIds = make(map[string]string)
	}
	p.recordIds[rr] = recordId

	return nil
}

func (p *dnsProvider) CleanUp(domain, token, keyAuth string) error {

	// 获取记录
	sub := strings.TrimSuffix(domain, "."+p.domain)
	rr := fmt.Sprintf("_acme-challenge.%s", sub)
	if sub == domain {
		rr = "_acme-challenge"
	}

	// 获取记录ID
	recordId := p.recordIds[rr]

	// 删除DNS记录
	return p.providerClient.DeleteDns(p.domain, recordId)
}

// RequestDomainCertificate 完成证书申请
func (c *certificate) RequestDomainCertificate(data *DomainCertificateRequest) error {

	// 创建数据库记录
	rrList := strings.Split(data.RR, ",")
	var fullDomains []string
	for _, rr := range rrList {
		fullDomains = append(fullDomains, fmt.Sprintf("%s.%s", rr, data.Domain))
	}
	crt := &model.DomainCertificate{
		Domain: strings.Join(fullDomains, " | "),
		Status: "pending",
	}
	if err := global.MySQLClient.Create(crt).Error; err != nil {
		return nil
	}

	// 获取域名服务商配置信息
	provider, err := dao.Domain.GetDomainServiceProviderForID(int(data.ProviderID))
	if err != nil {
		return err
	}

	// 创建请求客户端
	providerClient, err := Domain.GetCloudProviderClient(provider)
	if err != nil {
		return err
	}

	// 创建私钥
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}

	// 创建 AcmeUser
	acmeUser := AcmeUser{
		Email: *data.Email,
		key:   key,
	}

	conf := lego.NewConfig(&acmeUser)

	// 配置证书请求地址，测试环境为：LEDirectoryStaging，生产环境为：LEDirectoryProduction
	conf.CADirURL = lego.LEDirectoryStaging

	// 设置证书类型
	conf.Certificate.KeyType = certcrypto.RSA2048

	client, err := lego.NewClient(conf)
	if err != nil {
		return err
	}

	// 注册 ACME 账户
	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return err
	}
	acmeUser.Registration = reg
	acmeUser.GetRegistration()

	// 验证DNS
	customDnsProvider := dnsProvider{
		providerClient: providerClient,
		domain:         data.Domain,
	}
	err = client.Challenge.SetDNS01Provider(&customDnsProvider)
	if err != nil {
		return err
	}

	// 申请证书
	request := cert.ObtainRequest{
		Domains: fullDomains,
		Bundle:  true,
	}

	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		return err
	}

	// 将证书和私钥存储到数据库中
	crt.Certificate = string(certificates.Certificate)
	crt.PrivateKey = string(certificates.PrivateKey)
	crt.Status = "active"

	// 解析证书信息
	certInfo, err := parseCertificate(string(certificates.Certificate))
	if err != nil {
		return err
	}
	crt.StartAt = &certInfo.NotBefore
	crt.ExpirationAt = &certInfo.NotAfter
	return dao.Certificate.UpdateDomainCertificate(crt)
}

// UploadDomainCertificate 上传证书
func (c *certificate) UploadDomainCertificate(data *DomainCertificateCreate) (res *model.DomainCertificate, err error) {

	// 解析证书信息
	certInfo, err := parseCertificate(data.Certificate)
	if err != nil {
		return nil, err
	}

	// 解析私钥
	privateKey, err := parsePrivateKey(data.PrivateKey)
	if err != nil {
		return nil, err
	}

	// 检查证书和私钥是否匹配
	if err := verifyKeyPair(certInfo, privateKey); err != nil {
		return nil, err
	}

	// 判断证书是否过期
	if certInfo.NotAfter.Before(time.Now()) {
		return nil, errors.New("证书已过期")
	}

	// 获取绑定的域名
	boundDomains := getCertificateDomains(certInfo)

	crt := &model.DomainCertificate{
		Certificate:  data.Certificate,
		Domain:       strings.Join(boundDomains, " | "),
		ExpirationAt: &certInfo.NotAfter,
		PrivateKey:   data.PrivateKey,
		ServerType:   data.ServerType,
		StartAt:      &certInfo.NotBefore,
		Type:         data.Type,
		Status:       "active",
	}

	return dao.Certificate.UploadDomainCertificate(crt)
}

// DeleteDomainCertificate 删除证书
func (c *certificate) DeleteDomainCertificate(id int) (err error) {
	return dao.Certificate.DeleteCertificate(id)
}

// GetDomainCertificateList 获取证书列表
func (c *certificate) GetDomainCertificateList(name string, page, limit int) (*dao.DomainCertificateList, error) {
	return dao.Certificate.GetDomainCertificateList(name, page, limit)
}

// DownloadDomainCertificate 下载证书列表
func (c *certificate) DownloadDomainCertificate(id int) (data *bytes.Buffer, domainName string, err error) {

	// 获取证书
	crt, err := dao.Certificate.GetCertificateForID(id)
	if err != nil {
		return nil, "", err
	}

	// 临时文件名（直接内存使用）
	parts := strings.Split(crt.Domain, "|")
	baseName := strings.TrimSpace(parts[0])
	certFileName := fmt.Sprintf("%s.crt", baseName)
	keyFileName := fmt.Sprintf("%s.pem", baseName)

	// 创建内存 buffer
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// 将 cert 内容直接写入 zip
	if err := addContentToZip(zipWriter, certFileName, []byte(crt.Certificate)); err != nil {
		return nil, "", err
	}

	// 将 private key 内容直接写入 zip
	if err := addContentToZip(zipWriter, keyFileName, []byte(crt.PrivateKey)); err != nil {
		return nil, "", err
	}

	// 完成 zip
	_ = zipWriter.Close()

	// 4. 删除临时文件
	_ = os.Remove(certFileName)
	_ = os.Remove(keyFileName)

	return buf, strings.TrimSpace(parts[0]), nil
}

// CertificateExpiredNotice 证书过期通知
func (c *certificate) CertificateExpiredNotice() error {

	crts, err := dao.Certificate.GetExpiredCertificateList()
	if err != nil {
		return err
	}

	if len(crts) == 0 {
		logger.Info("检查证书状态正常.")
		return nil
	}

	// 生成HTML内容
	htmlBody := certificateExpiredNoticeHTML(crts)

	// 发送邮件函数
	return mail.Email.SendMsg([]string{"270142877@qq.com"}, nil, nil, "证书过期提醒", htmlBody, "html")
}

// certificateExpiredNoticeHTML 证书过期通知 HTML
func certificateExpiredNoticeHTML(certificates []*model.DomainCertificate) string {

	var (
		issuer = config.Conf.Settings["issuer"].(string)
		now    = time.Now()
	)

	return fmt.Sprintf(`
        <!DOCTYPE html>
        <html>
        <head>
            <meta charset="UTF-8">
            <title>证书过期提醒</title>
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
                    <h1>证书过期提醒</h1>
                    <div class="info">
                        以下证书即将过期或已过期，请及时处理以避免服务中断。
                    </div>
                    
                    <table>
                        <thead>
                            <tr>
                                <th>证书</th>
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
    `, generateCertificateRows(certificates, now), issuer)
}

// generateCertificateRows 生成证书表格行
func generateCertificateRows(certificates []*model.DomainCertificate, now time.Time) string {
	var rows strings.Builder

	for _, certificate := range certificates {

		// 确定域名状态
		status := "正常"
		statusClass := ""

		if certificate.ExpirationAt.Before(now) {
			status = "已过期"
			statusClass = "expired"
		} else if certificate.ExpirationAt.Before(now.Add(30 * 24 * time.Hour)) {
			status = "即将过期"
			statusClass = "expiring"
		}

		// 格式化过期时间
		expiredAt := certificate.ExpirationAt.Format("2006-01-02 15:04:05")

		rows.WriteString(fmt.Sprintf(`
            <tr>
                <td>%s</td>
                <td>%s</td>
                <td class="%s">%s</td>
            </tr>
        `, certificate.Domain, expiredAt, statusClass, status))
	}

	return rows.String()
}

// 直接写入内存数据到zip
func addContentToZip(zipWriter *zip.Writer, filename string, data []byte) error {
	writer, err := zipWriter.Create(filename)
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	return err
}

// parseCertificate 解析证书，提取相关信息
func parseCertificate(certPEM string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return nil, errors.New("无效的证书")
	}

	crt, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	return crt, nil
}

// parsePrivateKey 解析私钥
func parsePrivateKey(keyPEM string) (interface{}, error) {
	block, _ := pem.Decode([]byte(keyPEM))
	if block == nil {
		return nil, errors.New("无效的私钥")
	}

	// 私钥类型判断
	switch block.Type {

	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)

	case "EC PRIVATE KEY":
		return x509.ParseECPrivateKey(block.Bytes)

	case "PRIVATE KEY":
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}

		switch key.(type) {
		case *rsa.PrivateKey, *ecdsa.PrivateKey:
			return key, nil
		default:
			return nil, errors.New("不支持的私钥类型")
		}
	default:
		return nil, errors.New("未知的私钥类型")
	}
}

// getCertificateDomains 提取证书绑定的域名
func getCertificateDomains(cert *x509.Certificate) []string {
	var domains []string
	return append(domains, cert.DNSNames...)
}

// verifyKeyPair 验证证书和私钥是否匹配
func verifyKeyPair(cert *x509.Certificate, key interface{}) error {
	pubKey := cert.PublicKey

	switch pub := pubKey.(type) {
	case *rsa.PublicKey:
		priv, ok := key.(*rsa.PrivateKey)
		if !ok || priv.PublicKey.N.Cmp(pub.N) != 0 {
			return errors.New("RSA 公私钥不匹配")
		}
	case *ecdsa.PublicKey:
		priv, ok := key.(*ecdsa.PrivateKey)
		if !ok || priv.PublicKey.X.Cmp(pub.X) != 0 || priv.PublicKey.Y.Cmp(pub.Y) != 0 {
			return errors.New("ECDSA 公私钥不匹配")
		}
	default:
		return errors.New("不支持的密钥类型")
	}
	return nil
}
