package public_cloud

import "time"

// DomainList 域名列表
type DomainList struct {
	Name                    string     `json:"name"`
	RegistrationAt          *time.Time `json:"registration_at"`
	ExpirationAt            *time.Time `json:"expiration_at"`
	DomainServiceProviderID uint       `json:"domain_service_provider_id"`
}

// DnsList 域名DNS解析列表
type DnsList struct {
	Items []*DNS `json:"items"`
	Total int64  `json:"total"`
}

// DNS DNS记录
type DNS struct {
	RR       string `json:"rr"`
	Type     string `json:"type"`
	Value    string `json:"value"`
	TTL      int    `json:"ttl"`
	Status   string `json:"status"`
	CreateAt string `json:"create_at"`
	Remark   string `json:"remark"`
	RecordId string `json:"record_id"`
	Priority int    `json:"priority"`
	Weight   *int   `json:"weight"`
}
