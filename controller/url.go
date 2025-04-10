package controller

import (
	"github.com/gin-gonic/gin"
	"ops-api/dao"
	"ops-api/service"
	"strconv"
)

var UrlAddress urlAddress

type urlAddress struct{}

// AddUrl 新增Url
// @Summary 新增Url
// @Description Url监控相关
// @Tags Url监控相关
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param user body service.UrlAddressCreate true "Url信息"
// @Success 200 {string} json "{"code": 0, "msg": "创建成功", "data": nil}"
// @Router /api/v1/url [post]
func (u *urlAddress) AddUrl(c *gin.Context) {
	var url = &service.UrlAddressCreate{}

	if err := c.ShouldBind(url); err != nil {
		Response(c, 90400, err.Error())
		return
	}

	data, err := service.UrlAddress.AddUrl(url)
	if err != nil {
		Response(c, 90500, err.Error())
		return
	}

	CreateOrUpdateResponse(c, 0, "创建成功", data)
}

// DeleteUrl 删除Url
// @Summary 删除Url
// @Description Url监控相关
// @Tags Url监控相关
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path int true "ID"
// @Success 200 {string} json "{"code": 0, "msg": "删除成功"}"
// @Router /api/v1/url/{id} [delete]
func (u *urlAddress) DeleteUrl(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		Response(c, 90500, err.Error())
		return
	}

	if err := service.UrlAddress.DeleteUrl(id); err != nil {
		Response(c, 90500, err.Error())
		return
	}

	Response(c, 0, "删除成功")
}

// UpdateUrl 更新Url
// @Summary 更新Url
// @Description Url监控相关
// @Tags Url监控相关
// @Param Authorization header string true "Bearer 用户令牌"
// @Param task body dao.UrlAddressUpdate true "域名信息"
// @Success 200 {string} json "{"code": 0, "msg": "更新成功", "data": nil}"
// @Router /api/v1/url [put]
func (u *urlAddress) UpdateUrl(c *gin.Context) {
	var data = &dao.UrlAddressUpdate{}
	if err := c.ShouldBind(&data); err != nil {
		Response(c, 90400, err.Error())
		return
	}

	provider, err := service.UrlAddress.UpdateUrl(data)
	if err != nil {
		Response(c, 90500, err.Error())
		return
	}

	CreateOrUpdateResponse(c, 0, "更新成功", provider)
}

// GetUrlList 获取Url列表
// @Summary 获取Url列表
// @Description Url监控相关
// @Tags Url监控相关
// @Param Authorization header string true "Bearer 用户令牌"
// @Param page query int true "分页"
// @Param limit query int true "分页大小"
// @Param name query string false "Url信息"
// @Success 200 {string} json "{"code": 0, "data": []}"
// @Router /api/v1/urls [get]
func (u *urlAddress) GetUrlList(c *gin.Context) {

	params := new(struct {
		Name  string `form:"name"`
		Page  *int   `form:"page"`
		Limit *int   `form:"limit"`
	})

	if err := c.Bind(params); err != nil {
		Response(c, 90400, err.Error())
		return
	}

	data, err := service.UrlAddress.GetUrlList(params.Name, params.Page, params.Limit)
	if err != nil {
		Response(c, 90500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": data,
	})
}

// CertificateCheck 证书检查
// @Summary 证书检查
// @Description Url监控相关
// @Tags Url监控相关
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id query int true "ID"
// @Success 200 {string} json "{"code": 0, "msg": "检查完成"}"
// @Router /api/v1/url/check [post]
func (u *urlAddress) CertificateCheck(c *gin.Context) {

	params := new(struct {
		ID *uint `form:"id"`
	})

	if err := c.Bind(params); err != nil {
		Response(c, 90400, err.Error())
		return
	}

	if err := service.UrlAddress.CertificateCheck(params.ID); err != nil {
		Response(c, 90500, err.Error())
		return
	}

	Response(c, 0, "检查完成")
}
