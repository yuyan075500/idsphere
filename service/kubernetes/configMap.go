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

var ConfigMap configMap

type configMap struct{}

type ConfigMapList struct {
	Items *[]corev1.ConfigMap `json:"items"`
	Total int                 `json:"total"`
}

// List 获取ConfigMap列表
func (c *configMap) List(uuid, name, namespace string, page, limit int) (*ConfigMapList, error) {
	client := global.KubernetesClients.GetClient(uuid)
	if client == nil {
		return nil, fmt.Errorf("cluster %v not found", uuid)
	}

	configmaps, err := client.ClientSet.CoreV1().ConfigMaps(namespace).List(context.TODO(), metav1.ListOptions{})

	// 名称过滤
	var filtered []corev1.ConfigMap
	if name != "" {
		for _, item := range configmaps.Items {
			if strings.Contains(item.Name, name) {
				filtered = append(filtered, item)
			}
		}
	} else {
		filtered = configmaps.Items
	}

	// 分页
	res, err := utils.Paginate(filtered, page, limit)
	if err != nil {
		return nil, err
	}

	return &ConfigMapList{
		Items: res.(*[]corev1.ConfigMap),
		Total: len(filtered),
	}, nil
}
