package controller

import (
	"github.com/gin-gonic/gin"
	"ops-api/dao"
	"ops-api/service"
	"strconv"
)

var Domain domain

type domain struct{}

// AddDomainServiceProvider 新增域名服务商
// @Summary 新增域名服务商
// @Description 域名相关
// @Tags 域名管理
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param user body service.DomainServiceProviderCreate true "域名服务商信息"
// @Success 200 {string} json "{"code": 0, "msg": "创建成功", "data": nil}"
// @Router /api/v1/domain/provider [post]
func (d *domain) AddDomainServiceProvider(c *gin.Context) {
	var provider = &service.DomainServiceProviderCreate{}

	if err := c.ShouldBind(provider); err != nil {
		Response(c, 90400, err.Error())
		return
	}

	data, err := service.Domain.AddDomainServiceProvider(provider)
	if err != nil {
		Response(c, 90500, err.Error())
		return
	}

	CreateOrUpdateResponse(c, 0, "创建成功", data)
}

// DeleteDomainServiceProvider 删除域名服务商
// @Summary 删除域名服务商
// @Description 域名相关
// @Tags 域名管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path int true "域名服务商ID"
// @Success 200 {string} json "{"code": 0, "msg": "删除成功"}"
// @Router /api/v1/domain/provider/{id} [delete]
func (d *domain) DeleteDomainServiceProvider(c *gin.Context) {

	providerId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		Response(c, 90500, err.Error())
		return
	}

	if err := service.Domain.DeleteDomainServiceProvider(providerId); err != nil {
		Response(c, 90500, err.Error())
		return
	}

	Response(c, 0, "删除成功")
}

// UpdateDomainServiceProvider 更新域名服务商
// @Summary 更新域名服务商
// @Description 域名相关
// @Tags 域名管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param task body dao.ProviderUpdate true "服务商信息"
// @Success 200 {string} json "{"code": 0, "msg": "更新成功", "data": nil}"
// @Router /api/v1/domain/provider [put]
func (d *domain) UpdateDomainServiceProvider(c *gin.Context) {
	var data = &dao.ProviderUpdate{}
	if err := c.ShouldBind(&data); err != nil {
		Response(c, 90400, err.Error())
		return
	}

	provider, err := service.Domain.UpdateDomainServiceProviderList(data)
	if err != nil {
		Response(c, 90500, err.Error())
		return
	}

	CreateOrUpdateResponse(c, 0, "更新成功", provider)
}

// GetDomainServiceProviderList 获取域名服务商列表
// @Summary 获取域名服务商列表
// @Description 域名相关
// @Tags 域名管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Success 200 {string} json "{"code": 0, "data": []}"
// @Router /api/v1/domain/providers [get]
func (d *domain) GetDomainServiceProviderList(c *gin.Context) {

	data, err := service.Domain.GetDomainServiceProviderList()
	if err != nil {
		Response(c, 90500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": data,
	})
}

// AddDomain 新增域名
// @Summary 新增域名
// @Description 域名相关
// @Tags 域名管理
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param user body service.DomainCreate true "域名服务商信息"
// @Success 200 {string} json "{"code": 0, "msg": "创建成功", "data": nil}"
// @Router /api/v1/domain [post]
func (d *domain) AddDomain(c *gin.Context) {
	var domain = &service.DomainCreate{}

	if err := c.ShouldBind(domain); err != nil {
		Response(c, 90400, err.Error())
		return
	}

	data, err := service.Domain.AddDomain(domain)
	if err != nil {
		Response(c, 90500, err.Error())
		return
	}

	CreateOrUpdateResponse(c, 0, "创建成功", data)
}

// DeleteDomain 删除域名
// @Summary 删除域名
// @Description 域名相关
// @Tags 域名管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path int true "域名ID"
// @Success 200 {string} json "{"code": 0, "msg": "删除成功"}"
// @Router /api/v1/domain/{id} [delete]
func (d *domain) DeleteDomain(c *gin.Context) {

	providerId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		Response(c, 90500, err.Error())
		return
	}

	if err := service.Domain.DeleteDomain(providerId); err != nil {
		Response(c, 90500, err.Error())
		return
	}

	Response(c, 0, "删除成功")
}

// UpdateDomain 更新域名
// @Summary 更新域名
// @Description 域名相关
// @Tags 域名管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param task body dao.DomainUpdate true "域名信息"
// @Success 200 {string} json "{"code": 0, "msg": "更新成功", "data": nil}"
// @Router /api/v1/domain [put]
func (d *domain) UpdateDomain(c *gin.Context) {
	var data = &dao.DomainUpdate{}
	if err := c.ShouldBind(&data); err != nil {
		Response(c, 90400, err.Error())
		return
	}

	provider, err := service.Domain.UpdateDomain(data)
	if err != nil {
		Response(c, 90500, err.Error())
		return
	}

	CreateOrUpdateResponse(c, 0, "更新成功", provider)
}

// GetDomainList 获取域名列表
// @Summary 获取域名列表
// @Description 域名相关
// @Tags 域名管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param page query int true "分页"
// @Param limit query int true "分页大小"
// @Param name query string false "域名信息"
// @Success 200 {string} json "{"code": 0, "data": []}"
// @Router /api/v1/domains [get]
func (d *domain) GetDomainList(c *gin.Context) {

	params := new(struct {
		Name       string `form:"name"`
		ProviderId uint   `form:"provider_id"`
		Page       int    `form:"page" binding:"required"`
		Limit      int    `form:"limit" binding:"required"`
	})

	if err := c.Bind(params); err != nil {
		Response(c, 90400, err.Error())
		return
	}

	data, err := service.Domain.GetDomainList(params.Name, params.ProviderId, params.Page, params.Limit)
	if err != nil {
		Response(c, 90500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": data,
	})
}
