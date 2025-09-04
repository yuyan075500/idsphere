package kubernetes

import (
	"github.com/gin-gonic/gin"
	"ops-api/controller"
	dao "ops-api/dao/kubernetes"
	service "ops-api/service/kubernetes"
	"strconv"
)

var Cluster cluster

type cluster struct{}

// AddCluster 新增集群
// @Summary 新增集群
// @Description Kubernetes相关接口
// @Tags Kubernetes相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param user body service.CreateData true "集群信息"
// @Success 200 {string} json "{"code": 0, "msg": "创建成功", "data": nil}"
// @Router /api/v1/kubernetes/cluster [post]
func (cl *cluster) AddCluster(c *gin.Context) {
	var params = &service.CreateData{}

	if err := c.ShouldBindBodyWithJSON(params); err != nil {
		controller.Response(c, 90400, err.Error())
		return
	}

	data, err := service.Cluster.AddCluster(params)
	if err != nil {
		controller.Response(c, 90500, err.Error())
		return
	}

	controller.CreateOrUpdateResponse(c, 0, "创建成功", data)
}

// DeleteCluster 删除集群
// @Summary 删除集群
// @Description Kubernetes相关接口
// @Tags Kubernetes相关接口
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path int true "ID"
// @Success 200 {string} json "{"code": 0, "msg": "删除成功"}"
// @Router /api/v1/kubernetes/cluster/{id} [delete]
func (cl *cluster) DeleteCluster(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		controller.Response(c, 90500, err.Error())
		return
	}

	if err := service.Cluster.DeleteCluster(id); err != nil {
		controller.Response(c, 90500, err.Error())
		return
	}

	controller.Response(c, 0, "删除成功")
}

// UpdateCluster 修改集群信息
// @Summary 修改集群信息
// @Description Kubernetes相关接口
// @Tags Kubernetes相关接口
// @Param Authorization header string true "Bearer 用户令牌"
// @Param task body dao.UpdateData true "集群信息"
// @Success 200 {string} json "{"code": 0, "msg": "更新成功", "data": nil}"
// @Router /api/v1/kubernetes/cluster [put]
func (cl *cluster) UpdateCluster(c *gin.Context) {
	var data = &dao.UpdateData{}
	if err := c.ShouldBindBodyWithJSON(&data); err != nil {
		controller.Response(c, 90400, err.Error())
		return
	}

	provider, err := service.Cluster.UpdateCluster(data)
	if err != nil {
		controller.Response(c, 90500, err.Error())
		return
	}

	controller.CreateOrUpdateResponse(c, 0, "更新成功", provider)
}

// GetKubernetesList 获取集群列表
// @Summary 获取集群列表
// @Description Kubernetes相关接口
// @Tags Kubernetes相关接口
// @Param Authorization header string true "Bearer 用户令牌"
// @Param page query int true "分页"
// @Param limit query int true "分页大小"
// @Param name query string false "集群名称"
// @Success 200 {string} json "{"code": 0, "data": []}"
// @Router /api/v1/kubernetes/clusters [get]
func (cl *cluster) GetKubernetesList(c *gin.Context) {

	params := new(struct {
		Name  string `form:"name"`
		Page  int    `form:"page" binding:"required"`
		Limit int    `form:"limit" binding:"required"`
	})
	if err := c.ShouldBindQuery(params); err != nil {
		controller.Response(c, 90400, err.Error())
		return
	}

	data, err := service.Cluster.GetKubernetesList(params.Name, params.Page, params.Limit)

	if err != nil {
		controller.Response(c, 90500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": data,
	})
}
