package kubernetes

import (
	"bytes"
	"context"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/printers"
	"ops-api/kubernetes"
	"ops-api/utils"
	"strings"
)

var CronJob cronjob

type cronjob struct{}

type CronJobList struct {
	Items *[]batchv1.CronJob `json:"items"`
	Total int                `json:"total"`
}

// List 获取CronJob列表
func (c *cronjob) List(name, namespace string, page, limit int, client *kubernetes.ClientList) (*CronJobList, error) {
	deployments, err := client.ClientSet.BatchV1().CronJobs(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// 名称过滤
	var filtered []batchv1.CronJob
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
	res, err := utils.Paginate(filtered, page, limit)
	if err != nil {
		return nil, err
	}

	return &CronJobList{
		Items: res.(*[]batchv1.CronJob),
		Total: len(filtered),
	}, nil
}

// GetYAML 获取StatefulSet YAML配置
func (c *cronjob) GetYAML(name, namespace string, client *kubernetes.ClientList) (string, error) {
	data, err := client.ClientSet.BatchV1().CronJobs(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	// 设置GroupVersionKind
	data.GetObjectKind().SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "batch",
		Version: "v1",
		Kind:    "CronJob",
	})

	// 输出 YAML
	buf := new(bytes.Buffer)
	y := printers.YAMLPrinter{}
	if err := y.PrintObj(data, buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}
