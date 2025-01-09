package check

import (
	"errors"
	"fmt"
	"ops-api/config"
	"strings"
)

// PasswordCheck 检查密码复杂度，用于检查密码复杂度
func PasswordCheck(password string) error {

	var (
		passwordLength                            = config.Conf.Settings["passwordLength"].(int)
		passwordComplexity                        = config.Conf.Settings["passwordComplexity"].([]string)
		hasUpper, hasLower, hasNumber, hasSpecial bool
	)

	// 长度检验
	if len(password) < passwordLength {
		return errors.New(fmt.Sprintf("密码长度至少为 %d 位", passwordLength))
	}

	// 包含的类型校验
	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case strings.ContainsAny(string(char), "!@#$%^&*()-_=+[]{}|;:'\",.<>?/"):
			hasSpecial = true
		}
	}
	for _, rule := range passwordComplexity {
		switch rule {
		case "numbers":
			if !hasNumber {
				return errors.New("必须包含数字")
			}
		case "uppercase":
			if !hasUpper {
				return errors.New("必须包含大写字母")
			}
		case "lowercase":
			if !hasLower {
				return errors.New("必须包含小写字母")
			}
		case "specialCharacters":
			if !hasSpecial {
				return errors.New("必须包含特殊字符")
			}
		}
	}

	return nil
}
