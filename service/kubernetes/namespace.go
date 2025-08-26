package kubernetes

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"ops-api/global"
	"ops-api/kubernetes"
	"ops-api/utils"
	"strings"
)

var Namespace namespace

type namespace struct{}

type NamespaceList struct {
	Items *[]v1.Namespace `json:"items"`
	Total int             `json:"total"`
}

// List 获取命名空间列表
func (n *namespace) List(name string, page, limit int, client *kubernetes.ClientList) (*NamespaceList, error) {
	namespaces, err := client.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// 名称过滤
	var filtered []v1.Namespace
	if name != "" {
		for _, item := range namespaces.Items {
			if strings.Contains(item.Name, name) {
				filtered = append(filtered, item)
			}
		}
	} else {
		filtered = namespaces.Items
	}

	// 分页
	res, err := utils.Paginate(filtered, page, limit)
	if err != nil {
		return nil, err
	}

	return &NamespaceList{
		Items: res.(*[]v1.Namespace),
		Total: len(filtered),
	}, nil
}

// ListAll 获取命名空间列表（所有）
func (n *namespace) ListAll(uuid string) (*v1.NamespaceList, error) {
	client := global.KubernetesClients.GetClient(uuid)
	if client == nil {
		return nil, fmt.Errorf("cluster %v not found", uuid)
	}
	return client.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
}
