package kubernetes

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"ops-api/kubernetes"
	"ops-api/utils"
	"strings"
)

var Secret secret

type secret struct{}

type SecretList struct {
	Items *[]corev1.Secret `json:"items"`
	Total int              `json:"total"`
}

// List 获取Secret列表
func (s *secret) List(name, namespace string, page, limit int, client *kubernetes.ClientList) (*SecretList, error) {
	secrets, err := client.ClientSet.CoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// 名称过滤
	var filtered []corev1.Secret
	if name != "" {
		for _, item := range secrets.Items {
			if strings.Contains(item.Name, name) {
				filtered = append(filtered, item)
			}
		}
	} else {
		filtered = secrets.Items
	}

	// 分页
	res, err := utils.Paginate(filtered, page, limit)
	if err != nil {
		return nil, err
	}

	return &SecretList{
		Items: res.(*[]corev1.Secret),
		Total: len(filtered),
	}, nil
}
