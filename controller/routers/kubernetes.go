package routers

import (
	"github.com/gin-gonic/gin"
	controller "ops-api/controller/kubernetes"
)

func initKubernetesRouters(router *gin.Engine) {
	router.GET("/api/v1/clusters", controller.Cluster.GetKubernetesList)
	cluster := router.Group("/api/v1/cluster")
	{
		// 新增集群
		cluster.POST("", controller.Cluster.AddCluster)
		// 删除
		cluster.DELETE("/:id", controller.Cluster.DeleteCluster)
		// 修改
		cluster.PUT("", controller.Cluster.UpdateCluster)
	}

	// 节点相关
	router.GET("/api/v1/kubernetes/nodes", controller.Node.ListNodes)
	node := router.Group("/api/v1/kubernetes/node")
	{
		node.GET("/:name", controller.Node.GetYAML)
	}

	// 名称空间相关
	router.GET("/api/v1/kubernetes/namespaces", controller.Namespace.ListNamespaces)
	namespace := router.Group("/api/v1/kubernetes/namespace")
	{
		namespace.GET("/all", controller.Namespace.ListNamespacesAll)
		namespace.GET("/:name", controller.Namespace.GetYAML)
	}

	// Pod相关
	router.GET("/api/v1/kubernetes/pods", controller.Pod.ListPods)
	pod := router.Group("/api/v1/kubernetes/pod")
	{
		pod.GET("/:name", controller.Pod.GetYAML)
	}

	// Deployment相关
	router.GET("/api/v1/kubernetes/deployments", controller.Deployment.ListDeployments)

	// DaemonSet相关
	router.GET("/api/v1/kubernetes/daemonSets", controller.DaemonSet.ListDaemonSets)

	// StatefulSet相关
	router.GET("/api/v1/kubernetes/statefulSets", controller.StatefulSet.ListStatefulSets)

	// Job相关
	router.GET("/api/v1/kubernetes/jobs", controller.Job.ListJobs)

	// CronJob相关
	router.GET("/api/v1/kubernetes/cronJobs", controller.CronJob.ListCronJobs)

	// Service相关
	router.GET("/api/v1/kubernetes/services", controller.Svc.ListServices)

	// Endpoint相关
	router.GET("/api/v1/kubernetes/endpoints", controller.Endpoint.ListEndpoints)

	// Ingress相关
	router.GET("/api/v1/kubernetes/ingresses", controller.Ingress.ListIngresses)

	// PersistentVolume相关
	router.GET("/api/v1/kubernetes/persistentVolumes", controller.PersistentVolume.ListPersistentVolumes)

	// PersistentVolumeClaim相关
	router.GET("/api/v1/kubernetes/persistentVolumeClaims", controller.PersistentVolumeClaim.ListPersistentVolumeClaims)

	// StorageClass相关
	router.GET("/api/v1/kubernetes/storageClasses", controller.StorageClass.ListStorageClasses)

	// ConfigMap相关
	router.GET("/api/v1/kubernetes/configMaps", controller.ConfigMap.ListConfigMaps)

	// Secret相关
	router.GET("/api/v1/kubernetes/secrets", controller.Secret.ListSecrets)
}
