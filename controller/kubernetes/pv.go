package kubernetes

import (
	"github.com/gin-gonic/gin"
	"ops-api/controller"
	"ops-api/kubernetes"
	service "ops-api/service/kubernetes"
)

var PersistentVolume persistentVolume

type persistentVolume struct{}

// ListPersistentVolumes 获取PersistentVolume列表
func (p *persistentVolume) ListPersistentVolumes(c *gin.Context) {

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
	list, err := service.PersistentVolume.List(params.Name, params.Page, params.Limit, client)
	if err != nil {
		controller.Response(c, 90500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": list,
	})
}

// GetYAML 获取节点PersistentVolume YAML配置
func (p *persistentVolume) GetYAML(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		controller.Response(c, 90400, "名称不能为空")
		return
	}

	client := c.MustGet("kc").(*kubernetes.ClientList)
	strData, err := service.PersistentVolume.GetYAML(name, client)
	if err != nil {
		controller.Response(c, 90500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": strData,
	})
}
