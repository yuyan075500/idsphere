package kubernetes

import (
	"github.com/gin-gonic/gin"
	"ops-api/controller"
	"ops-api/kubernetes"
	svr "ops-api/service/kubernetes"
)

var Endpoint endpoint

type endpoint struct{}

// ListEndpoints 获取Endpoint列表
func (e *endpoint) ListEndpoints(c *gin.Context) {

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
	list, err := svr.Endpoint.List(params.Name, params.Namespace, params.Page, params.Limit, client)
	if err != nil {
		controller.Response(c, 90500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": list,
	})
}

// GetYAML 获取Endpoint YAML配置
func (e *endpoint) GetYAML(c *gin.Context) {
	uriParams := new(struct {
		Name string `uri:"name" binding:"required"`
	})
	if err := c.ShouldBindUri(uriParams); err != nil {
		controller.Response(c, 90400, err.Error())
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
	strData, err := svr.Endpoint.GetYAML(uriParams.Name, params.Namespace, client)
	if err != nil {
		controller.Response(c, 90500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": strData,
	})
}
