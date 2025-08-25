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

var Svc svc

type svc struct{}

type ServiceList struct {
	Items *[]corev1.Service `json:"items"`
	Total int               `json:"total"`
}

// List 获取Service列表
func (s *svc) List(uuid, name, namespace string, page, limit int) (*ServiceList, error) {
	client := global.KubernetesClients.GetClient(uuid)
	if client == nil {
		return nil, fmt.Errorf("cluster %v not found", uuid)
	}

	services, err := client.ClientSet.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})

	// 名称过滤
	var filtered []corev1.Service
	if name != "" {
		for _, item := range services.Items {
			if strings.Contains(item.Name, name) {
				filtered = append(filtered, item)
			}
		}
	} else {
		filtered = services.Items
	}

	// 分页
	res, err := utils.Paginate(filtered, page, limit)
	if err != nil {
		return nil, err
	}

	return &ServiceList{
		Items: res.(*[]corev1.Service),
		Total: len(filtered),
	}, nil
}
