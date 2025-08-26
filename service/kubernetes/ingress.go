package kubernetes

import (
	"context"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"ops-api/kubernetes"
	"ops-api/utils"
	"strings"
)

var Ingress ingress

type ingress struct{}

type IngressList struct {
	Items *[]networkingv1.Ingress `json:"items"`
	Total int                     `json:"total"`
}

// List 获取Ingress列表
func (i *ingress) List(name, namespace string, page, limit int, client *kubernetes.ClientList) (*IngressList, error) {
	ingresses, err := client.ClientSet.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// 名称过滤
	var filtered []networkingv1.Ingress
	if name != "" {
		for _, item := range ingresses.Items {
			if strings.Contains(item.Name, name) {
				filtered = append(filtered, item)
			}
		}
	} else {
		filtered = ingresses.Items
	}

	// 分页
	res, err := utils.Paginate(filtered, page, limit)
	if err != nil {
		return nil, err
	}

	return &IngressList{
		Items: res.(*[]networkingv1.Ingress),
		Total: len(filtered),
	}, nil
}
