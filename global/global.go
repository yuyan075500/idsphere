package global

import (
	"github.com/casbin/casbin/v2"
	"github.com/go-redis/redis"
	"github.com/minio/minio-go/v7"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"ops-api/kubernetes"
)

// 全局变量
var (
	MinioClient       *minio.Client
	RedisClient       *redis.Client
	MySQLClient       *gorm.DB
	CasBinServer      *casbin.Enforcer
	CornSchedule      *cron.Cron
	KubernetesClients *kubernetes.Clients
)
