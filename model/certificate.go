package model

import "time"

type DomainCertificate struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Certificate  string    `json:"certificate"`
	PrivateKey   string    `json:"private_key"`
	Domain       string    `json:"domain"`
	Type         uint      `json:"type"`
	ServerType   uint      `json:"server_type"`
	StartAt      time.Time `json:"start_at"`
	ExpirationAt time.Time `json:"expiration_at"`
}
