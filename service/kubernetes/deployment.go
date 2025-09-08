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

var Deployment deployment

type deployment struct{}

type DeploymentList struct {
	Items *[]appsv1.Deployment `json:"items"`
	Total int                  `json:"total"`
}

// DeploymentBatchDeleteStruct 批量删除
type DeploymentBatchDeleteStruct struct {
	Deployments []DeploymentItem `json:"deployments" binding:"required,min=1"`
}
type DeploymentItem struct {
	Name      string `json:"name" binding:"required"`
	Namespace string `json:"namespace" binding:"required"`
}

// CreateFromYAML 通过YAML内容创建Deployment
func (d *deployment) CreateFromYAML(yamlContent []byte, clients *kubernetes.ClientList) (*appsv1.Deployment, error) {
	var deploy appsv1.Deployment
	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(yamlContent), 1024)
	if err := decoder.Decode(&deploy); err != nil {
		return nil, err
	}

	// 设置默认命名空间
	if deploy.Namespace == "" {
		deploy.Namespace = metav1.NamespaceDefault
	}

	return clients.ClientSet.AppsV1().Deployments(deploy.Namespace).Create(context.TODO(), &deploy, metav1.CreateOptions{})
}

// BatchDeleteDeployment 批量删除 Deployment
func (d *deployment) BatchDeleteDeployment(deployments []DeploymentItem, client *kubernetes.ClientList) error {

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

// List 获取Deployment列表
func (d *deployment) List(name, namespace string, page, limit int, client *kubernetes.ClientList) (*DeploymentList, error) {
	deployments, err := client.ClientSet.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// 名称过滤
	var filtered []appsv1.Deployment
	if name != "" {
		for _, item := range deployments.Items {
			if strings.Contains(item.Name, name) {
				filtered = append(filtered, item)
			}
		}
	} else {
		filtered = deployments.Items
	}

	// 分页
	res, err := Paginate(filtered, page, limit)
	if err != nil {
		return nil, err
	}

	return &DeploymentList{
		Items: res.(*[]appsv1.Deployment),
		Total: len(filtered),
	}, nil
}

// GetYAML 获取Deployment YAML配置
func (d *deployment) GetYAML(name, namespace string, client *kubernetes.ClientList) (string, error) {
	data, err := client.ClientSet.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	// 设置GroupVersionKind
	data.GetObjectKind().SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "apps",
		Version: "v1",
		Kind:    "Deployment",
	})

	// 输出 YAML
	buf := new(bytes.Buffer)
	y := printers.YAMLPrinter{}
	if err := y.PrintObj(data, buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}
