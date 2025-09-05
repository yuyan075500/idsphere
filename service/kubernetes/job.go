package kubernetes

import (
	"bytes"
	"context"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/printers"
	"ops-api/kubernetes"
	"strings"
)

var Job job

type job struct{}

type JobList struct {
	Items *[]batchv1.Job `json:"items"`
	Total int            `json:"total"`
}

// List 获取Job列表
func (j *job) List(name, namespace string, page, limit int, client *kubernetes.ClientList) (*JobList, error) {
	jobs, err := client.ClientSet.BatchV1().Jobs(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// 名称过滤
	var filtered []batchv1.Job
	if name != "" {
		for _, item := range jobs.Items {
			if strings.Contains(item.Name, name) {
				filtered = append(filtered, item)
			}
		}
	} else {
		filtered = jobs.Items
	}

	// 分页
	res, err := Paginate(filtered, page, limit)
	if err != nil {
		return nil, err
	}

	return &JobList{
		Items: res.(*[]batchv1.Job),
		Total: len(filtered),
	}, nil
}

// GetYAML 获取Job YAML配置
func (j *job) GetYAML(name, namespace string, client *kubernetes.ClientList) (string, error) {
	data, err := client.ClientSet.BatchV1().Jobs(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	// 设置GroupVersionKind
	data.GetObjectKind().SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "batch",
		Version: "v1",
		Kind:    "Job",
	})

	// 输出 YAML
	buf := new(bytes.Buffer)
	y := printers.YAMLPrinter{}
	if err := y.PrintObj(data, buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}
