package model

import (
	"gorm.io/gorm"
	"time"
)

// SsoOAuthTicket OAuth2认证票据
type SsoOAuthTicket struct {
	*gorm.Model
	ExpiresAt   time.Time  `json:"expires_at"`
	Code        string     `json:"code"`
	RedirectURI string     `json:"redirect_uri"`
	ConsumedAt  *time.Time `json:"consumed_at"`
	UserID      uint       `json:"user_id"`
	Nonce       *string    `json:"nonce"`
}

func (*SsoOAuthTicket) TableName() (name string) {
	return "sso_oauth_ticket"
}

// SsoCASTicket CAS认证票据
type SsoCASTicket struct {
	*gorm.Model
	ExpiresAt  time.Time  `json:"expires_at"`
	Ticket     string     `json:"ticket"`
	Service    string     `json:"service"`
	ConsumedAt *time.Time `json:"consumed_at"`
	UserID     uint       `json:"user_id"`
}

func (*SsoCASTicket) TableName() (name string) {
	return "sso_cas_ticket"
}

// SsoNginxTicket Nginx认证票据
type SsoNginxTicket struct {
	*gorm.Model
	ExpiresAt time.Time `json:"expires_at"`
	Token     string    `json:"code"`
	UserID    uint      `json:"user_id"`
}

func (*SsoNginxTicket) TableName() (name string) {
	return "sso_nginx_ticket"
}
