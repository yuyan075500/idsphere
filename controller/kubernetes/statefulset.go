package kubernetes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ops-api/controller"
	"ops-api/kubernetes"
	service "ops-api/service/kubernetes"
)

var StatefulSet statefulSet

type statefulSet struct{}

// CreateFromYAML 创建 StatefulSet
func (s *statefulSet) CreateFromYAML(c *gin.Context) {
	yamlData, err := c.GetRawData()
	if err != nil {
		controller.Response(c, 90400, err.Error())
		return
	}

	client := c.MustGet("kc").(*kubernetes.ClientList)
	_, err = service.StatefulSet.CreateFromYAML(yamlData, client)
	if err != nil {
		controller.Response(c, 90500, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "创建成功",
	})
}

// BatchDeleteDeployment 批量删除 StatefulSet
func (s *statefulSet) BatchDeleteDeployment(c *gin.Context) {

	var params = &service.StatefulSetBatchDeleteStruct{}
	if err := c.ShouldBind(params); err != nil {
		controller.Response(c, 90400, err.Error())
		return
	}

	client := c.MustGet("kc").(*kubernetes.ClientList)
	err := service.StatefulSet.BatchDeleteDaemonSet(params.StatefulSets, client)
	if err != nil {
		controller.Response(c, 90500, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "删除成功",
	})
}

// UpdateFromYAML 更新 StatefulSet
func (s *statefulSet) UpdateFromYAML(c *gin.Context) {
	yamlData, err := c.GetRawData()
	if err != nil {
		controller.Response(c, 90400, err.Error())
		return
	}

	client := c.MustGet("kc").(*kubernetes.ClientList)
	_, err = service.StatefulSet.UpdateFromYAML(yamlData, client)
	if err != nil {
		controller.Response(c, 90500, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "更新成功",
	})
}

// ListStatefulSets 获取StatefulSet列表
func (s *statefulSet) ListStatefulSets(c *gin.Context) {

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
	list, err := service.StatefulSet.List(params.Name, params.Namespace, params.Page, params.Limit, client)
	if err != nil {
		controller.Response(c, 90500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": list,
	})
}

// GetYAML 获取StatefulSet YAML配置
func (s *statefulSet) GetYAML(c *gin.Context) {
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
	strData, err := service.StatefulSet.GetYAML(uriParams.Name, params.Namespace, client)
	if err != nil {
		controller.Response(c, 90500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": strData,
	})
}
