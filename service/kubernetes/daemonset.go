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

var DaemonSet daemonSet

type daemonSet struct{}

type DaemonSetList struct {
	Items *[]appsv1.DaemonSet `json:"items"`
	Total int                 `json:"total"`
}

// List 获取DaemonSet列表
func (d *daemonSet) List(uuid, name, namespace string, page, limit int) (*DaemonSetList, error) {
	client := global.KubernetesClients.GetClient(uuid)
	if client == nil {
		return nil, fmt.Errorf("cluster %v not found", uuid)
	}

	daemonSets, err := client.ClientSet.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})

	// 名称过滤
	var filtered []appsv1.DaemonSet
	if name != "" {
		for _, item := range daemonSets.Items {
			if strings.Contains(item.Name, name) {
				filtered = append(filtered, item)
			}
		}
	} else {
		filtered = daemonSets.Items
	}

	// 分页
	res, err := utils.Paginate(filtered, page, limit)
	if err != nil {
		return nil, err
	}

	return &DaemonSetList{
		Items: res.(*[]appsv1.DaemonSet),
		Total: len(filtered),
	}, nil
}
