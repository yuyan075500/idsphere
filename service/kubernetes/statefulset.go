package kubernetes

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/cli-runtime/pkg/printers"
	"ops-api/kubernetes"
	"strings"
)

var StatefulSet statefulSet

type statefulSet struct{}

type StatefulSetList struct {
	Items *[]appsv1.StatefulSet `json:"items"`
	Total int                   `json:"total"`
}

// StatefulSetBatchDeleteStruct 批量删除
type StatefulSetBatchDeleteStruct struct {
	StatefulSets []StatefulSetItem `json:"statefulSets" binding:"required,min=1"`
}
type StatefulSetItem struct {
	Name      string `json:"name" binding:"required"`
	Namespace string `json:"namespace" binding:"required"`
}

// CreateFromYAML 通过 YAML 内容创建 StatefulSet
func (s *statefulSet) CreateFromYAML(yamlContent []byte, clients *kubernetes.ClientList) (*appsv1.StatefulSet, error) {
	var manifest appsv1.StatefulSet
	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(yamlContent), 1024)
	if err := decoder.Decode(&manifest); err != nil {
		return nil, err
	}

	// 设置默认命名空间
	if manifest.Namespace == "" {
		manifest.Namespace = metav1.NamespaceDefault
	}

	return clients.ClientSet.AppsV1().StatefulSets(manifest.Namespace).Create(context.TODO(), &manifest, metav1.CreateOptions{})
}

// BatchDeleteDaemonSet 批量删除 StatefulSet
func (s *statefulSet) BatchDeleteDaemonSet(statefulSets []StatefulSetItem, client *kubernetes.ClientList) error {

	var failed []string

	// 遍历删除
	for _, value := range statefulSets {
		if err := client.ClientSet.AppsV1().StatefulSets(value.Namespace).Delete(
			context.TODO(),
			value.Name,
			metav1.DeleteOptions{},
		); err != nil {
			logger.Error(fmt.Printf("StatefulSet %s 删除失败: %v", value.Name, err))
			failed = append(failed, value.Name)
		}
	}

	if len(failed) > 0 {
		return errors.New("部分 Deployments 删除失败: " + strings.Join(failed, ", "))
	}

	return nil
}

// UpdateFromYAML 通过 YAML 内容修改已存在的 StatefulSet
func (s *statefulSet) UpdateFromYAML(yamlContent []byte, client *kubernetes.ClientList) (*appsv1.StatefulSet, error) {
	var manifest appsv1.StatefulSet
	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(yamlContent), 1024)
	if err := decoder.Decode(&manifest); err != nil {
		return nil, err
	}

	existing, err := client.ClientSet.AppsV1().StatefulSets(manifest.Namespace).Get(context.TODO(), manifest.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	manifest.ResourceVersion = existing.ResourceVersion

	return client.ClientSet.AppsV1().StatefulSets(manifest.Namespace).Update(context.TODO(), &manifest, metav1.UpdateOptions{})
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
	res, err := Paginate(filtered, page, limit)
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
