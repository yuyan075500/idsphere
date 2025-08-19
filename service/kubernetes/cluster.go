package kubernetes

import (
	"fmt"
	"github.com/wonderivan/logger"
	"golang.org/x/net/context"
	dao "ops-api/dao/kubernetes"
	"ops-api/global"
	model "ops-api/model/kubernetes"
	u "ops-api/utils"
)

var Cluster cluster

type cluster struct{}

// CreateData 新增集群结构体
type CreateData struct {
	Name       string `json:"name" binding:"required"`
	AuthType   int    `json:"auth_type" binding:"required"`
	Kubeconfig string `json:"kubeconfig"`
	UUID       string `json:"uuid"`
}

// AddCluster 新增集群
func (c *cluster) AddCluster(data *CreateData) (res *model.Cluster, err error) {

	d := &model.Cluster{
		Name:       data.Name,
		AuthType:   data.AuthType,
		Kubeconfig: data.Kubeconfig,
		UUID:       u.GenerateRandomString(16),
	}

	result, err := dao.Cluster.AddCluster(d)
	if err != nil {
		return nil, err
	}

	// 重新加载 kubernetes 客户端
	_ = global.KubernetesClients.Reload(global.MySQLClient)

	return result, nil
}

// DeleteCluster 删除集群
func (c *cluster) DeleteCluster(id int) error {

	err := dao.Cluster.DeleteCluster(id)
	if err != nil {
		return err
	}

	// 重新加载 kubernetes 客户端
	_ = global.KubernetesClients.Reload(global.MySQLClient)

	return nil
}

// UpdateCluster 修改集群信息
func (c *cluster) UpdateCluster(data *dao.UpdateData) (*model.Cluster, error) {

	result, err := dao.Cluster.UpdateCluster(data)
	if err != nil {
		return nil, err
	}

	// 重新加载 kubernetes 客户端
	_ = global.KubernetesClients.Reload(global.MySQLClient)

	return result, nil
}

// GetKubernetesList 获取集群列表
func (c *cluster) GetKubernetesList(name string, page, limit int) (data *dao.K8sList, err error) {
	data, err = dao.Cluster.GetKubernetesList(name, page, limit)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *cluster) GetClusterInfo(uuid string) {

	// 获取客户端
	client := global.KubernetesClients.GetClient(uuid)

	// 获取版本信息
	version, err := client.DiscoveryClient.ServerVersion()
	if err != nil {
		return
	}

	// 获取健康状态
	status, err := client.ClientSet.RESTClient().
		Get().
		AbsPath("/healthz").
		Do(context.TODO()).
		Raw()
	if err != nil {
		logger.Warn("获取集群健康状态失败")
	}
	fmt.Println(version)
	fmt.Println(string(status))
}
