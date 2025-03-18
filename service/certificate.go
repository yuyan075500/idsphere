package service

import (
	"archive/zip"
	"bytes"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"ops-api/dao"
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

	// 获取绑定的域名
	boundDomains := getCertificateDomains(certInfo)

	cert := &model.DomainCertificate{
		Certificate:  data.Certificate,
		Domain:       strings.Join(boundDomains, " | "),
		ExpirationAt: certInfo.NotAfter,
		PrivateKey:   data.PrivateKey,
		ServerType:   data.ServerType,
		StartAt:      certInfo.NotBefore,
		Type:         data.Type,
	}

	fmt.Println(cert)

	return dao.Certificate.UploadDomainCertificate(cert)
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
	cert, err := dao.Certificate.GetCertificateForID(id)
	if err != nil {
		return nil, "", err
	}

	// 临时文件名（直接内存使用）
	parts := strings.Split(cert.Domain, "|")
	baseName := strings.TrimSpace(parts[0])
	certFileName := fmt.Sprintf("%s.crt", baseName)
	keyFileName := fmt.Sprintf("%s.pem", baseName)

	// 创建内存 buffer
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// 将 cert 内容直接写入 zip
	if err := addContentToZip(zipWriter, certFileName, []byte(cert.Certificate)); err != nil {
		return nil, "", err
	}

	// 将 private key 内容直接写入 zip
	if err := addContentToZip(zipWriter, keyFileName, []byte(cert.PrivateKey)); err != nil {
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

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	return cert, nil
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
