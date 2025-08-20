package kubernetes

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"ops-api/global"
	"ops-api/utils"
	"strings"
)

var Pod pod

type pod struct{}

type PodList struct {
	Items *[]v1.Pod `json:"items"`
	Total int       `json:"total"`
}

// List 获取命名空间列表
func (p *pod) List(uuid, name, namespace string, page, limit int) (*PodList, error) {
	client := global.KubernetesClients.GetClient(uuid)
	if client == nil {
		return nil, fmt.Errorf("cluster %v not found", uuid)
	}

	pods, err := client.ClientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})

	// 名称过滤
	var filtered []v1.Pod
	if name != "" {
		for _, item := range pods.Items {
			if strings.Contains(item.Name, name) {
				filtered = append(filtered, item)
			}
		}
	} else {
		filtered = pods.Items
	}

	// 分页
	res, err := utils.Paginate(filtered, page, limit)
	if err != nil {
		return nil, err
	}

	return &PodList{
		Items: res.(*[]v1.Pod),
		Total: len(filtered),
	}, nil
}
