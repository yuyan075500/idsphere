package kubernetes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	service "ops-api/service/kubernetes"
)

var Secret secret

type secret struct{}

// ListSecrets 获取Secret列表
func (s *secret) ListSecrets(c *gin.Context) {

	params := new(struct {
		UUID      string `form:"uuid" binding:"required"`
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

	list, err := service.Secret.List(params.UUID, params.Name, params.Namespace, params.Page, params.Limit)
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
