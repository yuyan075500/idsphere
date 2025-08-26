package kubernetes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ops-api/kubernetes"
	service "ops-api/service/kubernetes"
)

var StorageClass storageClass

type storageClass struct{}

// ListStorageClasses 获取StorageClass列表
func (s *storageClass) ListStorageClasses(c *gin.Context) {

	params := new(struct {
		Name  string `form:"name"`
		Limit int    `form:"limit" binding:"required"`
		Page  int    `form:"page" binding:"required"`
	})
	if err := c.Bind(params); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	client := c.MustGet("kc").(*kubernetes.ClientList)
	list, err := service.StorageClass.List(params.Name, params.Page, params.Limit, client)
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
