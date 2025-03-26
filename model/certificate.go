package model

import "time"

type DomainCertificate struct {
	ID           uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	Certificate  string     `json:"certificate"`
	PrivateKey   string     `json:"private_key"`
	Domain       string     `json:"domain"`
	Type         uint       `json:"type" gorm:"default:1"`
	ServerType   uint       `json:"server_type" gorm:"default:1"`
	StartAt      *time.Time `json:"start_at,omitempty"`
	ExpirationAt *time.Time `json:"expiration_at,omitempty"`
	Status       string     `json:"status"` // active：可用（未过期）, expired：已过期, pending：申请中
}
