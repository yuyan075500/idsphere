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

var PersistentVolume persistentVolume

type persistentVolume struct{}

type PersistentVolumeList struct {
	Items *[]corev1.PersistentVolume `json:"items"`
	Total int                        `json:"total"`
}

// List 获取PersistentVolume列表
func (p *persistentVolume) List(uuid, name string, page, limit int) (*PersistentVolumeList, error) {
	client := global.KubernetesClients.GetClient(uuid)
	if client == nil {
		return nil, fmt.Errorf("cluster %v not found", uuid)
	}

	persistentVolumes, err := client.ClientSet.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})

	// 名称过滤
	var filtered []corev1.PersistentVolume
	if name != "" {
		for _, item := range persistentVolumes.Items {
			if strings.Contains(item.Name, name) {
				filtered = append(filtered, item)
			}
		}
	} else {
		filtered = persistentVolumes.Items
	}

	// 分页
	res, err := utils.Paginate(filtered, page, limit)
	if err != nil {
		return nil, err
	}

	return &PersistentVolumeList{
		Items: res.(*[]corev1.PersistentVolume),
		Total: len(filtered),
	}, nil
}
