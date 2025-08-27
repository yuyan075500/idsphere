package kubernetes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ops-api/kubernetes"
	service "ops-api/service/kubernetes"
)

var Namespace namespace

type namespace struct{}

// ListNamespaces 获取命名空间列表
func (n *namespace) ListNamespaces(c *gin.Context) {

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
	list, err := service.Namespace.List(params.Name, params.Page, params.Limit, client)
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

// ListNamespacesAll 获取命名空间列表（所有）
func (n *namespace) ListNamespacesAll(c *gin.Context) {

	params := new(struct {
		UUID string `form:"uuid" binding:"required"`
	})
	if err := c.Bind(params); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	list, err := service.Namespace.ListAll(params.UUID)
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

// GetYAML 获取节点YAML配置
func (n *namespace) GetYAML(c *gin.Context) {

	var (
		name   = c.Param("name")
		client = c.MustGet("kc").(*kubernetes.ClientList)
	)

	strData, err := service.Namespace.GetYAML(name, client)
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
