package kubernetes

import (
	"fmt"
	"github.com/wonderivan/logger"
	"gorm.io/gorm"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	model "ops-api/model/kubernetes"
)

// Clients 集群客户端
type Clients struct {
	clientMap map[string]*ClientList
}

// ClientList 集群客户端列表
type ClientList struct {
	ClientSet       *kubernetes.Clientset
	DiscoveryClient *discovery.DiscoveryClient
	DynamicClient   dynamic.Interface
	RestConfig      *rest.Config
}

// KubernetesInit 初始化 kubernetes 客户端
func (k *Clients) KubernetesInit(db *gorm.DB) error {

	var clusters []model.Cluster

	// 查找所有集群
	if err := db.Find(&clusters).Error; err != nil {
		return err
	}

	// 初始化 clientMap，并指定长度为集群的数量
	k.clientMap = make(map[string]*ClientList, len(clusters))

	for _, cluster := range clusters {

		// 将 kubeconfig 字符串转换成 rest.Config
		restConfig, err := clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Kubeconfig))
		if err != nil {
			logger.Error(fmt.Sprintf("解析集群 %s kubeconfig 失败: %v", cluster.Name, err))
			continue
		}

		// 创建 kubernetes clientset
		clientSet, err := kubernetes.NewForConfig(restConfig)
		if err != nil {
			logger.Error(fmt.Sprintf("Kubernetes集群（%s）clientset 失败: %v", cluster.Name, err))
			continue
		}

		discoveryClient, err := discovery.NewDiscoveryClientForConfig(restConfig)
		if err != nil {
			logger.Error(fmt.Sprintf("Kubernetes集群（%s）discoveryClient 失败: %v", cluster.Name, err))
			continue
		}

		dynamicClient, err := dynamic.NewForConfig(restConfig)
		if err != nil {
			logger.Error(fmt.Errorf("kubernetes集群（%s）dynamicClient 失败: %v", cluster.Name, err))
			continue
		}

		k.clientMap[cluster.UUID] = &ClientList{
			ClientSet:       clientSet,
			DiscoveryClient: discoveryClient,
			DynamicClient:   dynamicClient,
			RestConfig:      restConfig,
		}

		logger.Info(fmt.Sprintf("Kubernetes集群(%s)初始化成功.", cluster.Name))
	}

	return nil
}

// GetClient 获取 kubernetes 客户端
func (k *Clients) GetClient(uuid string) *ClientList {
	return k.clientMap[uuid]
}

// Reload 重新初始化 kubernetes 配置
func (k *Clients) Reload(db *gorm.DB) error {
	return k.KubernetesInit(db)
}
