package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"ops-api/config"
)

// Decrypt 字符串解密
func Decrypt(cipherText string) (string, error) {
	// 对Base64编码的字符串解码
	str, err := base64.RawURLEncoding.DecodeString(cipherText)

	// 读取私钥
	privateKeySrt := config.Conf.Settings["privateKey"].(string)

	// 私钥字符串转换为字节切片
	privateKeyData := []byte(privateKeySrt)

	block, _ := pem.Decode(privateKeyData)
	if block == nil || block.Type != "PRIVATE KEY" {
		return "", errors.New("invalid private key")
	}

	// 解析私钥
	privateKeyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// 私钥转换
	privateKey, ok := privateKeyInterface.(*rsa.PrivateKey)
	if !ok {
		return "", err
	}

	// 解密
	data, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, str)
	return string(data), err
}
