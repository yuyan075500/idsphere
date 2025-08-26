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

	// PersistentVolume相关
	pv := router.Group("/api/v1/kubernetes/persistentVolume")
	{
		pv.GET("/", controller.PersistentVolume.ListPersistentVolumes)
	}

	// PersistentVolumeClaim相关
	pvc := router.Group("/api/v1/kubernetes/persistentVolumeClaim")
	{
		pvc.GET("/", controller.PersistentVolumeClaim.ListPersistentVolumeClaims)
	}

	// StorageClass相关
	sc := router.Group("/api/v1/kubernetes/storageClass")
	{
		sc.GET("/", controller.StorageClass.ListStorageClasses)
	}

	// ConfigMap相关
	configmap := router.Group("/api/v1/kubernetes/configMap")
	{
		configmap.GET("/", controller.ConfigMap.ListConfigMaps)
	}

	// Secret相关
	secret := router.Group("/api/v1/kubernetes/secret")
	{
		secret.GET("/", controller.Secret.ListSecrets)
	}
}
