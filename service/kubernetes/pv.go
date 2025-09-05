package kubernetes

import (
	"bytes"
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/printers"
	"ops-api/kubernetes"
	"strings"
)

var PersistentVolume persistentVolume

type persistentVolume struct{}

type PersistentVolumeList struct {
	Items *[]corev1.PersistentVolume `json:"items"`
	Total int                        `json:"total"`
}

// List 获取PersistentVolume列表
func (p *persistentVolume) List(name string, page, limit int, client *kubernetes.ClientList) (*PersistentVolumeList, error) {

	persistentVolumes, err := client.ClientSet.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

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
	res, err := Paginate(filtered, page, limit)
	if err != nil {
		return nil, err
	}

	return &PersistentVolumeList{
		Items: res.(*[]corev1.PersistentVolume),
		Total: len(filtered),
	}, nil
}

// GetYAML 获取PersistentVolume YAML配置
func (p *persistentVolume) GetYAML(name string, client *kubernetes.ClientList) (string, error) {
	data, err := client.ClientSet.CoreV1().PersistentVolumes().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	// 设置GroupVersionKind
	data.GetObjectKind().SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "PersistentVolume",
	})

	// 输出 YAML
	buf := new(bytes.Buffer)
	y := printers.YAMLPrinter{}
	if err := y.PrintObj(data, buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}
