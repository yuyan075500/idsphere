package kubernetes

import (
	"context"
	"fmt"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"ops-api/global"
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
func (c *cronjob) List(uuid, name, namespace string, page, limit int) (*CronJobList, error) {
	client := global.KubernetesClients.GetClient(uuid)
	if client == nil {
		return nil, fmt.Errorf("cluster %v not found", uuid)
	}

	deployments, err := client.ClientSet.BatchV1().CronJobs(namespace).List(context.TODO(), metav1.ListOptions{})

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
