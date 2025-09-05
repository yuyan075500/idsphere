package kubernetes

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/printers"
	"ops-api/kubernetes"
	"sort"
	"strings"
)

var Pod pod

type pod struct{}

type PodList struct {
	Items *[]v1.Pod `json:"items"`
	Total int       `json:"total"`
}

// BatchDeleteStruct 批量删除
type BatchDeleteStruct struct {
	Pods  []Item `json:"pods" binding:"required,min=1"`
	Force *int64 `json:"force"`
}
type Item struct {
	Name      string `json:"name" binding:"required"`
	Namespace string `json:"namespace" binding:"required"`
}

// BatchDeletePod 批量删除 Pod
func (p *pod) BatchDeletePod(pods []Item, seconds *int64, client *kubernetes.ClientList) error {

	// 删除选项，当 seconds=0 时表示强制删除
	deleteOptions := metav1.DeleteOptions{
		GracePeriodSeconds: seconds,
	}

	var failed []string

	// 遍历删除
	for _, value := range pods {
		if err := client.ClientSet.CoreV1().Pods(value.Namespace).Delete(
			context.TODO(),
			value.Name,
			deleteOptions,
		); err != nil {
			logger.Error(fmt.Printf("Pod %s 删除失败: %v", value.Name, err))
			failed = append(failed, value.Name)
		}
	}

	if len(failed) > 0 {
		return errors.New("部分Pod删除失败: " + strings.Join(failed, ", "))
	}

	return nil
}

// List 获取 Pod 列表
func (p *pod) List(name, namespace string, page, limit int, client *kubernetes.ClientList) (*PodList, error) {
	pods, err := client.ClientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// 名称过滤
	var filtered []v1.Pod
	if name != "" {
		for _, item := range pods.Items {
			if strings.Contains(item.Name, name) {
				filtered = append(filtered, item)
			}
		}
	} else {
		filtered = pods.Items
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

	return &PodList{
		Items: res.(*[]v1.Pod),
		Total: len(filtered),
	}, nil
}

// GetYAML 获取 Pod YAML 配置
func (p *pod) GetYAML(name, namespace string, client *kubernetes.ClientList) (string, error) {
	data, err := client.ClientSet.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	// 设置GroupVersionKind
	data.GetObjectKind().SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Pod",
	})

	// 输出 YAML
	buf := new(bytes.Buffer)
	y := printers.YAMLPrinter{}
	if err := y.PrintObj(data, buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}
