package kubernetes

import (
	"github.com/gin-gonic/gin"
	"io"
	v1 "k8s.io/api/core/v1"
	"ops-api/kubernetes"
)

var PodLog podLog

type podLog struct{}

func (p *podLog) StreamPodLog(c *gin.Context, namespace, podName, containerName string, follow bool, lines int64, client *kubernetes.ClientList) error {

	req := client.ClientSet.CoreV1().Pods(namespace).GetLogs(podName, &v1.PodLogOptions{
		Container: containerName,
		Follow:    follow, // 实时追踪日志
		TailLines: &lines, // 指定最近多少行
	})

	// 建立流
	stream, err := req.Stream(c)
	if err != nil {
		return err
	}
	defer func() {
		_ = stream.Close()
	}()

	// 推送日志
	c.Stream(func(w io.Writer) bool {
		select {
		case <-c.Request.Context().Done():
			return false
		default:
			buf := make([]byte, 2000)
			n, err := stream.Read(buf)
			if err != nil {
				if err == io.EOF {
					return false
				}
				return false
			}
			c.SSEvent("message", string(buf[:n]))
			return true
		}
	})
	return nil
}
