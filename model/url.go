package model

import "time"

type DomainCertificateMonitor struct {
	ID           uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name         string     `json:"name"`
	Domain       string     `json:"domain"`
	Port         uint       `json:"port"`
	IPAddress    string     `json:"ip_address"`
	Status       *int       `json:"status"` // 0: 正常，1: 检查异常，2: 已过期，nil：未检查
	ExpirationAt *time.Time `json:"expiration_at,omitempty"`
	LastCheckAt  *time.Time `json:"last_check_at,omitempty"`
}
