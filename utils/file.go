package utils

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"ops-api/config"
)

// LoadPublicKey 读取公钥
func LoadPublicKey() (interface{}, error) {

	// 读取公钥文件
	publicKeySrt := config.Conf.Settings["publicKey"].(string)

	// 公钥字符串转换为字节切片
	publicKeyData := []byte(publicKeySrt)

	// 解析PEM块
	block, _ := pem.Decode(publicKeyData)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("invalid public key")
	}

	// 解析公钥
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return pubKey, nil
}
