package kubernetes

import (
	"bytes"
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/printers"
	"ops-api/kubernetes"
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
	res, err := Paginate(filtered, page, limit)
	if err != nil {
		return nil, err
	}

	return &NodeList{
		Items: res.(*[]v1.Node),
		Total: len(filtered),
	}, nil
}

// GetYAML 获取节点YAML配置
func (n *node) GetYAML(name string, client *kubernetes.ClientList) (string, error) {
	data, err := client.ClientSet.CoreV1().Nodes().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	// 设置GroupVersionKind
	data.GetObjectKind().SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Node",
	})

	// 输出 YAML
	buf := new(bytes.Buffer)
	y := printers.YAMLPrinter{}
	if err := y.PrintObj(data, buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}
