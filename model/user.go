package model

import (
	"github.com/wonderivan/logger"
	"gorm.io/gorm"
	"ops-api/utils"
	"time"
)

// AuthUser 用户信息表

type AuthUser struct {
	gorm.Model
	Name              string       `json:"name"`
	WwId              *string      `json:"ww_id" gorm:"unique"`    // 企业微信用户ID
	CtyunId           *string      `json:"ctyun_id" gorm:"unique"` // 天翼云用户ID
	Username          string       `json:"username" gorm:"unique"`
	Avatar            *string      `json:"avatar"`
	Password          string       `json:"password"`
	PhoneNumber       string       `json:"phone_number"`
	IsActive          bool         `json:"is_active"`
	Email             string       `json:"email"`
	LastLoginAt       *time.Time   `json:"last_login_at"`
	MFACode           *string      `json:"mfa_code"`
	PasswordExpiredAt *time.Time   `json:"password_expired_at"`
	UserFrom          string       `json:"user_from" gorm:"default:本地"`
	Groups            []*AuthGroup `json:"groups" gorm:"many2many:auth_user_groups"`
	Accounts          []*Account   `gorm:"many2many:account_users"`
}

func (*AuthUser) TableName() (name string) {
	return "auth_user"
}

// BeforeCreate 新用户创建时对密码字段加密，仅创建用户时候调用
func (u *AuthUser) BeforeCreate(tx *gorm.DB) (err error) {
	cipherText, err := utils.Encrypt(u.Password)
	if err != nil {
		return err
	}
	u.Password = cipherText
	return nil
}

// CheckPassword 检查用户密码是否正确
func (u *AuthUser) CheckPassword(password string) bool {

	// 对数据库中的密码解密
	str, err := utils.Decrypt(u.Password)
	if err != nil {
		logger.Error(err.Error())
		return false
	}

	// 判断密码是否相等
	if str != password {
		return false
	}

	return true
}
