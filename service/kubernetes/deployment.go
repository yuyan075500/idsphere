package kubernetes

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"ops-api/global"
	"ops-api/utils"
	"strings"
)

var Deployment deployment

type deployment struct{}

type DeploymentList struct {
	Items *[]appsv1.Deployment `json:"items"`
	Total int                  `json:"total"`
}

// List 获取Deployment列表
func (d *deployment) List(uuid, name, namespace string, page, limit int) (*DeploymentList, error) {
	client := global.KubernetesClients.GetClient(uuid)
	if client == nil {
		return nil, fmt.Errorf("cluster %v not found", uuid)
	}

	deployments, err := client.ClientSet.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})

	// 名称过滤
	var filtered []appsv1.Deployment
	if name != "" {
		for _, item := range deployments.Items {
			if strings.Contains(item.Name, name) {
				filtered = append(filtered, item)
			}
		}
	} else {
		filtered = deployments.Items
	}

	// 分页
	res, err := utils.Paginate(filtered, page, limit)
	if err != nil {
		return nil, err
	}

	return &DeploymentList{
		Items: res.(*[]appsv1.Deployment),
		Total: len(filtered),
	}, nil
}
