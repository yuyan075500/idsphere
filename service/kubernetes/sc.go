package kubernetes

import (
	"bytes"
	"context"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/printers"
	"ops-api/kubernetes"
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
func (s *storageClass) List(name string, page, limit int, client *kubernetes.ClientList) (*StorageClassList, error) {

	storageClasses, err := client.ClientSet.StorageV1().StorageClasses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

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

// GetYAML 获取StorageClass YAML配置
func (s *storageClass) GetYAML(name string, client *kubernetes.ClientList) (string, error) {
	data, err := client.ClientSet.StorageV1().StorageClasses().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	// 设置GroupVersionKind
	data.GetObjectKind().SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "storage.k8s.io",
		Version: "v1",
		Kind:    "StorageClass",
	})

	// 输出 YAML
	buf := new(bytes.Buffer)
	y := printers.YAMLPrinter{}
	if err := y.PrintObj(data, buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}
