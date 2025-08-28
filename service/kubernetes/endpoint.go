package kubernetes

import (
	"bytes"
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/printers"
	"ops-api/kubernetes"
	"ops-api/utils"
	"strings"
)

var Endpoint endpoint

type endpoint struct{}

type EndpointList struct {
	Items *[]corev1.Endpoints `json:"items"`
	Total int                 `json:"total"`
}

// List 获取Endpoint列表
func (e *endpoint) List(name, namespace string, page, limit int, client *kubernetes.ClientList) (*EndpointList, error) {
	endpoints, err := client.ClientSet.CoreV1().Endpoints(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// 名称过滤
	var filtered []corev1.Endpoints
	if name != "" {
		for _, item := range endpoints.Items {
			if strings.Contains(item.Name, name) {
				filtered = append(filtered, item)
			}
		}
	} else {
		filtered = endpoints.Items
	}

	// 分页
	res, err := utils.Paginate(filtered, page, limit)
	if err != nil {
		return nil, err
	}

	return &EndpointList{
		Items: res.(*[]corev1.Endpoints),
		Total: len(filtered),
	}, nil
}

// GetYAML 获取Endpoints YAML配置
func (e *endpoint) GetYAML(name, namespace string, client *kubernetes.ClientList) (string, error) {
	data, err := client.ClientSet.CoreV1().Endpoints(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	// 设置GroupVersionKind
	data.GetObjectKind().SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Endpoints",
	})

	// 输出 YAML
	buf := new(bytes.Buffer)
	y := printers.YAMLPrinter{}
	if err := y.PrintObj(data, buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}
