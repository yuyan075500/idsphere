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

var StatefulSet statefulSet

type statefulSet struct{}

type StatefulSetList struct {
	Items *[]appsv1.StatefulSet `json:"items"`
	Total int                   `json:"total"`
}

// List 获取StatefulSet列表
func (s *statefulSet) List(uuid, name, namespace string, page, limit int) (*StatefulSetList, error) {
	client := global.KubernetesClients.GetClient(uuid)
	if client == nil {
		return nil, fmt.Errorf("cluster %v not found", uuid)
	}

	statefulSets, err := client.ClientSet.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})

	// 名称过滤
	var filtered []appsv1.StatefulSet
	if name != "" {
		for _, item := range statefulSets.Items {
			if strings.Contains(item.Name, name) {
				filtered = append(filtered, item)
			}
		}
	} else {
		filtered = statefulSets.Items
	}

	// 分页
	res, err := utils.Paginate(filtered, page, limit)
	if err != nil {
		return nil, err
	}

	return &StatefulSetList{
		Items: res.(*[]appsv1.StatefulSet),
		Total: len(filtered),
	}, nil
}
