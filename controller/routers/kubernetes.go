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
		namespace.POST("", controller.Namespace.CreateNamespaces)
		namespace.DELETE("/:name", controller.Namespace.DeleteNamespaces)
		namespace.PUT("", controller.Namespace.UpdateFromYAML)
		namespace.GET("/all", controller.Namespace.ListNamespacesAll)
		namespace.GET("/:name", controller.Namespace.GetYAML)
	}

	// Pod相关
	router.GET("/api/v1/kubernetes/pods", controller.Pod.ListPods)
	pod := router.Group("/api/v1/kubernetes/pod")
	{
		pod.GET("/:name", controller.Pod.GetYAML)
		pod.GET("/terminal", controller.PodTerminal.Init)
		pod.GET("/log", controller.PodLog.GetPodLogs)
	}

	// Deployment相关
	router.GET("/api/v1/kubernetes/deployments", controller.Deployment.ListDeployments)
	deployment := router.Group("/api/v1/kubernetes/deployment")
	{
		deployment.GET("/:name", controller.Deployment.GetYAML)
	}

	// DaemonSet相关
	router.GET("/api/v1/kubernetes/daemonSets", controller.DaemonSet.ListDaemonSets)
	daemonSet := router.Group("/api/v1/kubernetes/daemonSet")
	{
		daemonSet.GET("/:name", controller.DaemonSet.GetYAML)
	}

	// StatefulSet相关
	router.GET("/api/v1/kubernetes/statefulSets", controller.StatefulSet.ListStatefulSets)
	statefulSet := router.Group("/api/v1/kubernetes/statefulSet")
	{
		statefulSet.GET("/:name", controller.StatefulSet.GetYAML)
	}

	// Job相关
	router.GET("/api/v1/kubernetes/jobs", controller.Job.ListJobs)
	job := router.Group("/api/v1/kubernetes/job")
	{
		job.GET("/:name", controller.Job.GetYAML)
	}

	// CronJob相关
	router.GET("/api/v1/kubernetes/cronJobs", controller.CronJob.ListCronJobs)
	cronJob := router.Group("/api/v1/kubernetes/cronJob")
	{
		cronJob.GET("/:name", controller.CronJob.GetYAML)
	}

	// Service相关
	router.GET("/api/v1/kubernetes/services", controller.Svc.ListServices)
	service := router.Group("/api/v1/kubernetes/service")
	{
		service.GET("/:name", controller.Svc.GetYAML)
	}

	// Endpoint相关
	router.GET("/api/v1/kubernetes/endpoints", controller.Endpoint.ListEndpoints)
	endpoint := router.Group("/api/v1/kubernetes/endpoint")
	{
		endpoint.GET("/:name", controller.Endpoint.GetYAML)
	}

	// Ingress相关
	router.GET("/api/v1/kubernetes/ingresses", controller.Ingress.ListIngresses)
	ingress := router.Group("/api/v1/kubernetes/ingress")
	{
		ingress.GET("/:name", controller.Ingress.GetYAML)
	}

	// PersistentVolume相关
	router.GET("/api/v1/kubernetes/persistentVolumes", controller.PersistentVolume.ListPersistentVolumes)
	persistentVolume := router.Group("/api/v1/kubernetes/persistentVolume")
	{
		persistentVolume.GET("/:name", controller.PersistentVolume.GetYAML)
	}

	// PersistentVolumeClaim相关
	router.GET("/api/v1/kubernetes/persistentVolumeClaims", controller.PersistentVolumeClaim.ListPersistentVolumeClaims)
	persistentVolumeClaim := router.Group("/api/v1/kubernetes/persistentVolumeClaim")
	{
		persistentVolumeClaim.GET("/:name", controller.PersistentVolumeClaim.GetYAML)
	}

	// StorageClass相关
	router.GET("/api/v1/kubernetes/storageClasses", controller.StorageClass.ListStorageClasses)
	storageClass := router.Group("/api/v1/kubernetes/storageClass")
	{
		storageClass.GET("/:name", controller.StorageClass.GetYAML)
	}

	// ConfigMap相关
	router.GET("/api/v1/kubernetes/configMaps", controller.ConfigMap.ListConfigMaps)
	configMap := router.Group("/api/v1/kubernetes/configMap")
	{
		configMap.GET("/:name", controller.ConfigMap.GetYAML)
	}

	// Secret相关
	router.GET("/api/v1/kubernetes/secrets", controller.Secret.ListSecrets)
	secret := router.Group("/api/v1/kubernetes/secret")
	{
		secret.GET("/:name", controller.Secret.GetYAML)
	}
}
