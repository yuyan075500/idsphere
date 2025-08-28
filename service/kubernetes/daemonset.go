package kubernetes

import (
	"bytes"
	"context"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/printers"
	"ops-api/kubernetes"
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
func (d *daemonSet) List(name, namespace string, page, limit int, client *kubernetes.ClientList) (*DaemonSetList, error) {
	daemonSets, err := client.ClientSet.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

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

// GetYAML 获取DaemonSet YAML配置
func (d *daemonSet) GetYAML(name, namespace string, client *kubernetes.ClientList) (string, error) {
	data, err := client.ClientSet.AppsV1().DaemonSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	// 设置GroupVersionKind
	data.GetObjectKind().SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "apps",
		Version: "v1",
		Kind:    "DaemonSet",
	})

	// 输出 YAML
	buf := new(bytes.Buffer)
	y := printers.YAMLPrinter{}
	if err := y.PrintObj(data, buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}
