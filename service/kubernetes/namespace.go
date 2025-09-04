package kubernetes

import (
	"bytes"
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/cli-runtime/pkg/printers"
	"ops-api/kubernetes"
	"ops-api/utils"
	"strings"
)

var Namespace namespace

type namespace struct{}

type NamespaceList struct {
	Items *[]v1.Namespace `json:"items"`
	Total int             `json:"total"`
}

// Create 创建命名空间
func (n *namespace) Create(namespaceName, description string, client *kubernetes.ClientList) (*v1.Namespace, error) {
	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespaceName,
			Annotations: map[string]string{
				"description": description,
			},
		},
	}

	return client.ClientSet.CoreV1().Namespaces().Create(
		context.TODO(),
		ns,
		metav1.CreateOptions{},
	)
}

// Delete 删除指定命名空间
func (n *namespace) Delete(namespaceName string, client *kubernetes.ClientList) error {
	return client.ClientSet.CoreV1().Namespaces().Delete(
		context.TODO(),
		namespaceName,
		metav1.DeleteOptions{},
	)
}

// UpdateFromYAML 通过YAML内容修改已存在的Namespace
func (n *namespace) UpdateFromYAML(yamlContent string, client *kubernetes.ClientList) (*v1.Namespace, error) {
	// 结构体绑定
	var ns v1.Namespace
	decoder := yaml.NewYAMLOrJSONDecoder(strings.NewReader(yamlContent), 1024)
	if err := decoder.Decode(&ns); err != nil {
		return nil, err
	}

	// 获取 ResourceVersion
	existingNs, err := client.ClientSet.CoreV1().Namespaces().Get(
		context.TODO(),
		ns.Name,
		metav1.GetOptions{},
	)
	if err != nil {
		return nil, err
	}

	// 赋值 ResourceVersion，避免并发修改冲突
	ns.ResourceVersion = existingNs.ResourceVersion

	// 执行更新
	updatedNs, err := client.ClientSet.CoreV1().Namespaces().Update(
		context.TODO(),
		&ns,
		metav1.UpdateOptions{},
	)
	if err != nil {
		return nil, err
	}

	return updatedNs, nil
}

// List 获取命名空间列表
func (n *namespace) List(name string, page, limit int, client *kubernetes.ClientList) (*NamespaceList, error) {
	namespaces, err := client.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// 名称过滤
	var filtered []v1.Namespace
	if name != "" {
		for _, item := range namespaces.Items {
			if strings.Contains(item.Name, name) {
				filtered = append(filtered, item)
			}
		}
	} else {
		filtered = namespaces.Items
	}

	// 分页
	res, err := utils.Paginate(filtered, page, limit)
	if err != nil {
		return nil, err
	}

	return &NamespaceList{
		Items: res.(*[]v1.Namespace),
		Total: len(filtered),
	}, nil
}

// ListAll 获取命名空间列表（所有）
func (n *namespace) ListAll(client *kubernetes.ClientList) (*v1.NamespaceList, error) {
	return client.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
}

// GetYAML 获取Namespace YAML配置
func (n *namespace) GetYAML(name string, client *kubernetes.ClientList) (string, error) {
	data, err := client.ClientSet.CoreV1().Namespaces().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	// 设置GroupVersionKind
	data.GetObjectKind().SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Namespace",
	})

	// 输出 YAML
	buf := new(bytes.Buffer)
	y := printers.YAMLPrinter{}
	if err := y.PrintObj(data, buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}
