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
		namespace.GET("/all", controller.Namespace.ListNamespacesAll)
	}

	// Pod相关
	pod := router.Group("/api/v1/kubernetes/pod")
	{
		pod.GET("/", controller.Pod.ListPods)
	}

	// Deployment相关
	deployment := router.Group("/api/v1/kubernetes/deployment")
	{
		deployment.GET("/", controller.Deployment.ListDeployments)
	}

	// DaemonSet相关
	daemonSet := router.Group("/api/v1/kubernetes/daemonSet")
	{
		daemonSet.GET("/", controller.DaemonSet.ListDaemonSets)
	}

	// StatefulSet相关
	statefulSet := router.Group("/api/v1/kubernetes/statefulSet")
	{
		statefulSet.GET("/", controller.StatefulSet.ListStatefulSets)
	}

	// Job相关
	job := router.Group("/api/v1/kubernetes/job")
	{
		job.GET("/", controller.Job.ListJobs)
	}

	// CronJob相关
	cronJob := router.Group("/api/v1/kubernetes/cronJob")
	{
		cronJob.GET("/", controller.CronJob.ListCronJobs)
	}

	// Service相关
	service := router.Group("/api/v1/kubernetes/service")
	{
		service.GET("/", controller.Svc.ListServices)
	}

	// Endpoint相关
	endpoint := router.Group("/api/v1/kubernetes/endpoint")
	{
		endpoint.GET("/", controller.Endpoint.ListEndpoints)
	}

	// Ingress相关
	ingress := router.Group("/api/v1/kubernetes/ingress")
	{
		ingress.GET("/", controller.Ingress.ListIngresses)
	}
}
