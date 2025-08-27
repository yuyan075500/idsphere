package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ops-api/kubernetes"
	"strings"
)

func KubernetesClientInit(clients *kubernetes.Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 判断路径是否以 /api/v1/kubernetes/ 开头
		if strings.HasPrefix(c.Request.URL.Path, "/api/v1/kubernetes/") {
			// 获取 Kubernetes 集群 UUID
			uuid := c.GetHeader("X-Kubernetes-Cluster-Uuid")
			if uuid == "" {
				c.JSON(http.StatusOK, gin.H{
					"code": 90500,
					"msg":  "无效的 UUID",
				})
				return
			} else {
				// 获取 Kubernetes 集群客户端
				client := clients.GetClient(uuid)
				if client == nil {
					c.JSON(http.StatusOK, gin.H{
						"code": 90500,
						"msg":  "客户端初始化失败",
					})
					return
				}
				c.Set("kc", client)
			}
		}
		c.Next()
	}
}
