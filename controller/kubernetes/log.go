package kubernetes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ops-api/kubernetes"
	service "ops-api/service/kubernetes"
)

var PodLog podLog

type podLog struct{}

func (p *podLog) GetPodLogs(c *gin.Context) {
	params := new(struct {
		UUID          string `form:"uuid" binding:"required"`
		PodName       string `form:"podName" binding:"required"`
		ContainerName string `form:"containerName" binding:"required"`
		Namespace     string `form:"namespace" binding:"required"`
		Line          int64  `form:"line,default=50"`
		Follow        bool   `form:"follow"`
	})

	if err := c.ShouldBindQuery(params); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	client := c.MustGet("kc").(*kubernetes.ClientList)
	if err := service.PodLog.StreamPodLog(c, params.Namespace, params.PodName, params.ContainerName, params.Follow, params.Line, client); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
	}
}
