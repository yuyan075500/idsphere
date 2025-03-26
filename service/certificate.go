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
	"ops-api/dao"
	"ops-api/global"
	"ops-api/model"
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
	rr             string
	domain         string
	RecordId       string
}

func (p *dnsProvider) Present(domain, token, keyAuth string) error {
	// 获取DNS记录信息
	info := dns01.GetChallengeInfo(domain, keyAuth)

	// 组装DNS记录名称
	rr := fmt.Sprintf("_acme-challenge.%s", p.rr)

	// 添加DNS记录
	recordId, err := p.providerClient.AddDns(p.domain, "TXT", rr, info.Value, "DNS-01挑战", 600, nil, 0)
	if err != nil {
		return err
	}

	p.RecordId = recordId

	return nil
}

func (p *dnsProvider) CleanUp(domain, token, keyAuth string) error {
	// 删除DNS记录
	return p.providerClient.DeleteDns(p.domain, p.RecordId)
}

// RequestDomainCertificate 完成证书申请
func (c *certificate) RequestDomainCertificate(data *DomainCertificateRequest) error {

	// 创建数据库记录
	crt := &model.DomainCertificate{
		Domain: fmt.Sprintf("%s.%s", data.RR, data.Domain),
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

	config := lego.NewConfig(&acmeUser)

	// 配置证书请求地址，测试环境为：LEDirectoryStaging，生产环境为：LEDirectoryProduction
	config.CADirURL = lego.LEDirectoryStaging

	// 设置证书类型
	config.Certificate.KeyType = certcrypto.RSA2048

	client, err := lego.NewClient(config)
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

	// 验证DNS，设置验证超时时间为15分钟
	customDnsProvider := dnsProvider{
		providerClient: providerClient,
		domain:         data.Domain,
		rr:             data.RR,
	}
	err = client.Challenge.SetDNS01Provider(&customDnsProvider)
	if err != nil {
		return err
	}

	// 申请证书
	request := cert.ObtainRequest{
		Domains: []string{fmt.Sprintf("%s.%s", data.RR, data.Domain)},
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
