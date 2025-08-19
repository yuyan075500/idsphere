package routers

import (
	"github.com/gin-gonic/gin"
	controller "ops-api/controller/kubernetes"
)

func initKubernetesRouters(router *gin.Engine) {
	cluster := router.Group("/api/v1/kubernetes")
	{
		// 获取集群列表
		cluster.GET("/clusters", controller.Cluster.GetKubernetesList)
		// 新增集群
		cluster.POST("/cluster", controller.Cluster.AddCluster)
		// 删除
		cluster.DELETE("/cluster/:id", controller.Cluster.DeleteCluster)
		// 修改
		cluster.PUT("/cluster", controller.Cluster.UpdateCluster)

		// 获取集群基本信息
		cluster.GET("/cluster/info", controller.Cluster.GetKubernetesInfo)
	}

	// 节点相关
	node := router.Group("/api/v1/kubernetes/node")
	{
		node.GET("/", controller.Node.ListNodes)
	}

	// 名称空间相关
	namespace := router.Group("/api/v1/kubernetes/namespace")
	{
		namespace.GET("/", controller.Namespace.ListNamespaces)
	}
}
