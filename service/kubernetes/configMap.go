package kubernetes

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"ops-api/kubernetes"
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
func (c *configMap) List(name, namespace string, page, limit int, client *kubernetes.ClientList) (*ConfigMapList, error) {
	configmaps, err := client.ClientSet.CoreV1().ConfigMaps(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

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
