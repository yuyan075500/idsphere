package kubernetes

import (
	"ops-api/global"
)

import (
	"ops-api/model/kubernetes"
)

var Cluster cluster

type cluster struct{}

// K8sList 集群列表
type K8sList struct {
	Items []*kubernetes.Cluster `json:"items"`
	Total int64                 `json:"total"`
}

// UpdateData 集群基本信息
type UpdateData struct {
	ID         uint   `json:"id" binding:"required"`
	Name       string `json:"name" binding:"required"`
	AuthType   int    `json:"auth_type" binding:"required"`
	Kubeconfig string `json:"kubeconfig"`
}

// AddCluster 新增集群
func (c *cluster) AddCluster(data *kubernetes.Cluster) (res *kubernetes.Cluster, err error) {
	if err := global.MySQLClient.Create(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

// DeleteCluster 删除集群
func (c *cluster) DeleteCluster(id int) (err error) {
	if err := global.MySQLClient.Where("id = ?", id).Unscoped().Delete(&kubernetes.Cluster{}).Error; err != nil {
		return err
	}
	return nil
}

// UpdateCluster 修改集群信息
func (c *cluster) UpdateCluster(data *UpdateData) (*kubernetes.Cluster, error) {

	res := &kubernetes.Cluster{}

	if err := global.MySQLClient.Model(res).Select("*").Where("id = ?", data.ID).Updates(data).Error; err != nil {
		return nil, err
	}

	// 查询更新后的账号信息并返回
	if err := global.MySQLClient.First(res, data.ID).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// GetKubernetesList 获取集群列表
func (c *cluster) GetKubernetesList(name string, page, limit int) (data *K8sList, err error) {

	// 定义数据的起始位置
	startSet := (page - 1) * limit

	// 定义返回的内容
	var (
		items []*kubernetes.Cluster
		total int64
	)

	// 获取菜单列表
	tx := global.MySQLClient.Model(&kubernetes.Cluster{}).
		Where("name like ?", "%"+name+"%").
		Omit("kubeconfig").
		Count(&total).
		Limit(limit).
		Offset(startSet).
		Find(&items)
	if tx.Error != nil {
		return nil, err
	}

	return &K8sList{
		Items: items,
		Total: total,
	}, nil
}
