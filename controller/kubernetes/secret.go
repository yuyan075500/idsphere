package kubernetes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ops-api/kubernetes"
	service "ops-api/service/kubernetes"
)

var Secret secret

type secret struct{}

// ListSecrets 获取Secret列表
func (s *secret) ListSecrets(c *gin.Context) {

	params := new(struct {
		Namespace string `form:"namespace"`
		Name      string `form:"name"`
		Limit     int    `form:"limit" binding:"required"`
		Page      int    `form:"page" binding:"required"`
	})
	if err := c.Bind(params); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	client := c.MustGet("kc").(*kubernetes.ClientList)
	list, err := service.Secret.List(params.Name, params.Namespace, params.Page, params.Limit, client)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": list,
	})
}

// GetYAML 获取Secret YAML配置
func (s *secret) GetYAML(c *gin.Context) {

	var (
		name   = c.Param("name")
		client = c.MustGet("kc").(*kubernetes.ClientList)
	)

	params := new(struct {
		Namespace string `form:"namespace"`
	})
	if err := c.Bind(params); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	strData, err := service.Secret.GetYAML(name, params.Namespace, client)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": strData,
	})
}
