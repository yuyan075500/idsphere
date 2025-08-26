package kubernetes

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"ops-api/kubernetes"
	"ops-api/utils"
	"strings"
)

var PersistentVolumeClaim persistentVolumeClaim

type persistentVolumeClaim struct{}

type PersistentVolumeClaimList struct {
	Items *[]corev1.PersistentVolumeClaim `json:"items"`
	Total int                             `json:"total"`
}

// List 获取PersistentVolumeClaim列表
func (p *persistentVolumeClaim) List(name, namespace string, page, limit int, client *kubernetes.ClientList) (*PersistentVolumeClaimList, error) {
	persistentVolumeClaims, err := client.ClientSet.CoreV1().PersistentVolumeClaims(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// 名称过滤
	var filtered []corev1.PersistentVolumeClaim
	if name != "" {
		for _, item := range persistentVolumeClaims.Items {
			if strings.Contains(item.Name, name) {
				filtered = append(filtered, item)
			}
		}
	} else {
		filtered = persistentVolumeClaims.Items
	}

	// 分页
	res, err := utils.Paginate(filtered, page, limit)
	if err != nil {
		return nil, err
	}

	return &PersistentVolumeClaimList{
		Items: res.(*[]corev1.PersistentVolumeClaim),
		Total: len(filtered),
	}, nil
}
