package kubernetes

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"ops-api/global"
	"ops-api/utils"
	"strings"
)

var Endpoint endpoint

type endpoint struct{}

type EndpointList struct {
	Items *[]corev1.Endpoints `json:"items"`
	Total int                 `json:"total"`
}

// List 获取Endpoint列表
func (e *endpoint) List(uuid, name, namespace string, page, limit int) (*EndpointList, error) {
	client := global.KubernetesClients.GetClient(uuid)
	if client == nil {
		return nil, fmt.Errorf("cluster %v not found", uuid)
	}

	endpoints, err := client.ClientSet.CoreV1().Endpoints(namespace).List(context.TODO(), metav1.ListOptions{})

	// 名称过滤
	var filtered []corev1.Endpoints
	if name != "" {
		for _, item := range endpoints.Items {
			if strings.Contains(item.Name, name) {
				filtered = append(filtered, item)
			}
		}
	} else {
		filtered = endpoints.Items
	}

	// 分页
	res, err := utils.Paginate(filtered, page, limit)
	if err != nil {
		return nil, err
	}

	return &EndpointList{
		Items: res.(*[]corev1.Endpoints),
		Total: len(filtered),
	}, nil
}
