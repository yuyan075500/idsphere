package kubernetes

type Cluster struct {
	ID           uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	UUID         string `json:"uuid" gorm:"unique"` // 集群唯一标识
	Name         string `json:"name" gorm:"unique"` // 集群版本
	AuthType     int    `json:"auth_type"`          // 认证类型
	Kubeconfig   string `json:"kubeconfig"`         // 认证文件
	Version      string `json:"version"`            // 集群版本
	HealthStatus string `json:"health_status"`      // 健康状态
}

func (*Cluster) TableName() (name string) {
	return "kubernetes_cluster"
}
