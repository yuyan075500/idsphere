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

var StatefulSet statefulSet

type statefulSet struct{}

type StatefulSetList struct {
	Items *[]appsv1.StatefulSet `json:"items"`
	Total int                   `json:"total"`
}

// List 获取StatefulSet列表
func (s *statefulSet) List(name, namespace string, page, limit int, client *kubernetes.ClientList) (*StatefulSetList, error) {
	statefulSets, err := client.ClientSet.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

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

// GetYAML 获取StatefulSet YAML配置
func (s *statefulSet) GetYAML(name, namespace string, client *kubernetes.ClientList) (string, error) {
	data, err := client.ClientSet.AppsV1().StatefulSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	// 设置GroupVersionKind
	data.GetObjectKind().SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "apps",
		Version: "v1",
		Kind:    "StatefulSet",
	})

	// 输出 YAML
	buf := new(bytes.Buffer)
	y := printers.YAMLPrinter{}
	if err := y.PrintObj(data, buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}
