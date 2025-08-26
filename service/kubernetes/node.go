package kubernetes

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"ops-api/kubernetes"
	"ops-api/utils"
	"strings"
)

var Node node

type node struct{}

type NodeList struct {
	Items *[]v1.Node `json:"items"`
	Total int        `json:"total"`
}

func (n *node) List(name string, page, limit int, client *kubernetes.ClientList) (*NodeList, error) {

	nodes, err := client.ClientSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// 名称过滤
	var filtered []v1.Node
	if name != "" {
		// 匹配名称和IP地址
		for _, item := range nodes.Items {
			if strings.Contains(item.Name, name) {
				filtered = append(filtered, item)
				continue
			}
			for _, addr := range item.Status.Addresses {
				if strings.Contains(addr.Address, name) {
					filtered = append(filtered, item)
					break
				}
			}
		}

	} else {
		filtered = nodes.Items
	}

	// 分页
	res, err := utils.Paginate(filtered, page, limit)
	if err != nil {
		return nil, err
	}

	return &NodeList{
		Items: res.(*[]v1.Node),
		Total: len(filtered),
	}, nil
}
