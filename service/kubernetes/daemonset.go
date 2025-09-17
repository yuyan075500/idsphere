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
	"sort"
	"strings"
)

var DaemonSet daemonSet

type daemonSet struct{}

type DaemonSetList struct {
	Items *[]appsv1.DaemonSet `json:"items"`
	Total int                 `json:"total"`
}

// DaemonSetBatchDeleteStruct 批量删除
type DaemonSetBatchDeleteStruct struct {
	DaemonSets []DaemonSetItem `json:"deployments" binding:"required,min=1"`
}
type DaemonSetItem struct {
	Name      string `json:"name" binding:"required"`
	Namespace string `json:"namespace" binding:"required"`
}

// CreateFromYAML 通过 YAML 内容创建 DaemonSet
func (d *daemonSet) CreateFromYAML(yamlContent []byte, clients *kubernetes.ClientList) (*appsv1.DaemonSet, error) {
	var manifest appsv1.DaemonSet
	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(yamlContent), 1024)
	if err := decoder.Decode(&manifest); err != nil {
		return nil, err
	}

	// 设置默认命名空间
	if manifest.Namespace == "" {
		manifest.Namespace = metav1.NamespaceDefault
	}

	return clients.ClientSet.AppsV1().DaemonSets(manifest.Namespace).Create(context.TODO(), &manifest, metav1.CreateOptions{})
}

// BatchDeleteDaemonSet 批量删除 DaemonSet
func (d *daemonSet) BatchDeleteDaemonSet(deployments []DaemonSetItem, client *kubernetes.ClientList) error {

	var failed []string

	// 遍历删除
	for _, value := range deployments {
		if err := client.ClientSet.AppsV1().Deployments(value.Namespace).Delete(
			context.TODO(),
			value.Name,
			metav1.DeleteOptions{},
		); err != nil {
			logger.Error(fmt.Printf("Deployment %s 删除失败: %v", value.Name, err))
			failed = append(failed, value.Name)
		}
	}

	if len(failed) > 0 {
		return errors.New("部分 Deployments 删除失败: " + strings.Join(failed, ", "))
	}

	return nil
}

// UpdateFromYAML 通过 YAML 内容修改已存在的 DaemonSet
func (d *daemonSet) UpdateFromYAML(yamlContent []byte, client *kubernetes.ClientList) (*appsv1.DaemonSet, error) {
	var manifest appsv1.DaemonSet
	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(yamlContent), 1024)
	if err := decoder.Decode(&manifest); err != nil {
		return nil, err
	}

	existing, err := client.ClientSet.AppsV1().DaemonSets(manifest.Namespace).Get(context.TODO(), manifest.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	manifest.ResourceVersion = existing.ResourceVersion

	return client.ClientSet.AppsV1().DaemonSets(manifest.Namespace).Update(context.TODO(), &manifest, metav1.UpdateOptions{})
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

	// 按创建时间排序
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].CreationTimestamp.After(filtered[j].CreationTimestamp.Time)
	})

	// 分页
	res, err := Paginate(filtered, page, limit)
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
