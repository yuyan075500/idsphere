package kubernetes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ops-api/controller"
	"ops-api/kubernetes"
	service "ops-api/service/kubernetes"
)

var Namespace namespace

type namespace struct{}

// CreateNamespaces 创建命名空间
func (n *namespace) CreateNamespaces(c *gin.Context) {

	params := new(struct {
		Name string `form:"name" binding:"required"`
	})
	if err := c.ShouldBind(params); err != nil {
		controller.Response(c, 90400, err.Error())
		return
	}

	client := c.MustGet("kc").(*kubernetes.ClientList)
	list, err := service.Namespace.Create(params.Name, client)
	if err != nil {
		controller.Response(c, 90500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": list,
	})
}

// DeleteNamespaces 删除命名空间
func (n *namespace) DeleteNamespaces(c *gin.Context) {
	namespaceName := c.Param("name")
	if namespaceName == "" {
		controller.Response(c, 90400, "命名空间名称不能为空")
		return
	}

	client := c.MustGet("kc").(*kubernetes.ClientList)
	err := service.Namespace.Delete(namespaceName, client)
	if err != nil {
		controller.Response(c, 90500, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "删除成功",
	})
}

// ListNamespaces 获取命名空间列表
func (n *namespace) ListNamespaces(c *gin.Context) {

	params := new(struct {
		Name  string `form:"name"`
		Limit int    `form:"limit" binding:"required"`
		Page  int    `form:"page" binding:"required"`
	})
	if err := c.Bind(params); err != nil {
		controller.Response(c, 90400, err.Error())
		return
	}

	client := c.MustGet("kc").(*kubernetes.ClientList)
	list, err := service.Namespace.List(params.Name, params.Page, params.Limit, client)
	if err != nil {
		controller.Response(c, 90500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": list,
	})
}

// ListNamespacesAll 获取命名空间列表（所有）
func (n *namespace) ListNamespacesAll(c *gin.Context) {

	client := c.MustGet("kc").(*kubernetes.ClientList)
	list, err := service.Namespace.ListAll(client)
	if err != nil {
		controller.Response(c, 90500, err.Error())
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
		controller.Response(c, 90500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": strData,
	})
}
