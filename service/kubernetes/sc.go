package kubernetes

import (
	"context"
	"fmt"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"ops-api/global"
	"ops-api/utils"
	"strings"
)

var StorageClass storageClass

type storageClass struct{}

type StorageClassList struct {
	Items *[]storagev1.StorageClass `json:"items"`
	Total int                       `json:"total"`
}

// List 获取StorageClass列表
func (s *storageClass) List(uuid, name string, page, limit int) (*StorageClassList, error) {
	client := global.KubernetesClients.GetClient(uuid)
	if client == nil {
		return nil, fmt.Errorf("cluster %v not found", uuid)
	}

	storageClasses, err := client.ClientSet.StorageV1().StorageClasses().List(context.TODO(), metav1.ListOptions{})

	// 名称过滤
	var filtered []storagev1.StorageClass
	if name != "" {
		for _, item := range storageClasses.Items {
			if strings.Contains(item.Name, name) {
				filtered = append(filtered, item)
			}
		}
	} else {
		filtered = storageClasses.Items
	}

	// 分页
	res, err := utils.Paginate(filtered, page, limit)
	if err != nil {
		return nil, err
	}

	return &StorageClassList{
		Items: res.(*[]storagev1.StorageClass),
		Total: len(filtered),
	}, nil
}
