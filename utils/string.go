package utils

import (
	"errors"
	"math/rand"
	"net"
	"strings"
	"time"
)

// Contains 查询字符串在一个列表中是否存在
func Contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// MapToJson Map转json
//func MapToJson(data interface{}) string {
//	byteStr, _ := json.Marshal(data)
//	return string(byteStr)
//}

// GenerateRandomNumber 生成6位非零开头的随机数字
func GenerateRandomNumber() int {
	firstDigit := rand.Intn(9) + 1
	otherDigits := rand.Intn(100000)
	result := firstDigit*100000 + otherDigits
	return result
}

// GenerateRandomString 生成随机字符串
func GenerateRandomString(n int) string {

	// 指定随机字符串包含的字符集
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// GetSubdomain 获取二级域名
func GetSubdomain(host string) (string, error) {

	// 检查是否为IP地址或localhost
	if net.ParseIP(host) != nil || host == "localhost" {
		return host, nil
	}

	// 拆分主机名
	parts := strings.Split(host, ".")

	// 检查是否是有效的域名
	if len(parts) < 2 {
		return "", errors.New("无效的域名")
	}

	// 返回最后两部分作为二级域名
	return strings.Join(parts[len(parts)-2:], "."), nil
}

// ParseTime 解析日期字符串为 *time.Time
func ParseTime(dateStr string) *time.Time {
	if dateStr == "" {
		return nil
	}
	parsedTime, err := time.Parse("2006-01-02 15:04:05", dateStr)
	if err != nil {
		return nil
	}
	return &parsedTime
}
