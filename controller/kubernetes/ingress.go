package kubernetes

import (
	"github.com/gin-gonic/gin"
	"ops-api/controller"
	"ops-api/kubernetes"
	svr "ops-api/service/kubernetes"
)

var Ingress ingresses

type ingresses struct{}

// ListIngresses 获取Ingress列表
func (i *ingresses) ListIngresses(c *gin.Context) {

	params := new(struct {
		Namespace string `form:"namespace"`
		Name      string `form:"name"`
		Limit     int    `form:"limit" binding:"required"`
		Page      int    `form:"page" binding:"required"`
	})
	if err := c.ShouldBindQuery(params); err != nil {
		controller.Response(c, 90400, err.Error())
		return
	}

	client := c.MustGet("kc").(*kubernetes.ClientList)
	list, err := svr.Ingress.List(params.Name, params.Namespace, params.Page, params.Limit, client)
	if err != nil {
		controller.Response(c, 90500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": list,
	})
}

// GetYAML 获取Ingress YAML配置
func (i *ingresses) GetYAML(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		controller.Response(c, 90400, "名称不能为空")
		return
	}

	params := new(struct {
		Namespace string `form:"namespace"`
	})
	if err := c.ShouldBindQuery(params); err != nil {
		controller.Response(c, 90400, err.Error())
		return
	}

	client := c.MustGet("kc").(*kubernetes.ClientList)
	strData, err := svr.Ingress.GetYAML(name, params.Namespace, client)
	if err != nil {
		controller.Response(c, 90500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": strData,
	})
}
