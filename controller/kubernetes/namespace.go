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
		Name        string `form:"name" binding:"required"`
		Description string `form:"description"`
	})
	if err := c.ShouldBindBodyWithJSON(params); err != nil {
		controller.Response(c, 90400, err.Error())
		return
	}

	client := c.MustGet("kc").(*kubernetes.ClientList)
	list, err := service.Namespace.Create(params.Name, params.Description, client)
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
	params := new(struct {
		Name string `uri:"name" binding:"required"`
	})
	if err := c.ShouldBindUri(params); err != nil {
		controller.Response(c, 90400, err.Error())
		return
	}

	client := c.MustGet("kc").(*kubernetes.ClientList)
	err := service.Namespace.Delete(params.Name, client)
	if err != nil {
		controller.Response(c, 90500, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "删除成功",
	})
}

// UpdateFromYAML 更新名称空间
func (n *namespace) UpdateFromYAML(c *gin.Context) {
	yamlData, err := c.GetRawData()
	if err != nil {
		controller.Response(c, 90400, err.Error())
		return
	}

	client := c.MustGet("kc").(*kubernetes.ClientList)
	_, err = service.Namespace.UpdateFromYAML(string(yamlData), client)
	if err != nil {
		controller.Response(c, 90500, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "更新成功",
	})
}

// ListNamespaces 获取命名空间列表
func (n *namespace) ListNamespaces(c *gin.Context) {

	params := new(struct {
		Name  string `form:"name"`
		Limit int    `form:"limit" binding:"required"`
		Page  int    `form:"page" binding:"required"`
	})
	if err := c.ShouldBindQuery(params); err != nil {
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

	uriParams := new(struct {
		Name string `uri:"name" binding:"required"`
	})
	if err := c.ShouldBindUri(uriParams); err != nil {
		controller.Response(c, 90400, err.Error())
		return
	}

	client := c.MustGet("kc").(*kubernetes.ClientList)
	strData, err := service.Namespace.GetYAML(uriParams.Name, client)
	if err != nil {
		controller.Response(c, 90500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": strData,
	})
}
