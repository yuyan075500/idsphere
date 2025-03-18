package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"ops-api/service"
	"strconv"
)

var Certificate certificate

type certificate struct{}

// UploadDomainCertificate 上传域名证书
// @Summary 新增域名证书
// @Description 域名证书相关
// @Tags 域名证书相关
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param user body service.DomainCertificateCreate true "域名证书信息"
// @Success 200 {string} json "{"code": 0, "msg": "创建成功", "data": nil}"
// @Router /api/v1/certificate [post]
func (cert *certificate) UploadDomainCertificate(c *gin.Context) {
	var provider = &service.DomainCertificateCreate{}

	if err := c.ShouldBind(provider); err != nil {
		Response(c, 90400, err.Error())
		return
	}

	data, err := service.Certificate.UploadDomainCertificate(provider)
	if err != nil {
		Response(c, 90500, err.Error())
		return
	}

	CreateOrUpdateResponse(c, 0, "创建成功", data)
}

// DeleteDomainCertificate 删除域名证书
// @Summary 删除域名证书
// @Description 域名证书相关
// @Tags 域名证书相关
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path int true "域名证书ID"
// @Success 200 {string} json "{"code": 0, "msg": "删除成功"}"
// @Router /api/v1/certificate/{id} [delete]
func (cert *certificate) DeleteDomainCertificate(c *gin.Context) {

	certId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		Response(c, 90500, err.Error())
		return
	}

	if err := service.Certificate.DeleteDomainCertificate(certId); err != nil {
		Response(c, 90500, err.Error())
		return
	}

	Response(c, 0, "删除成功")
}

// GetDomainCertificateList 获取域名证书列表
// @Summary 获取域名证书列表
// @Description 域名证书相关
// @Tags 域名证书相关
// @Param Authorization header string true "Bearer 用户令牌"
// @Param page query int true "分页"
// @Param limit query int true "分页大小"
// @Param name query string false "域名信息"
// @Success 200 {string} json "{"code": 0, "data": []}"
// @Router /api/v1/certificates [get]
func (cert *certificate) GetDomainCertificateList(c *gin.Context) {

	params := new(struct {
		Name  string `form:"name"`
		Page  int    `form:"page" binding:"required"`
		Limit int    `form:"limit" binding:"required"`
	})

	if err := c.Bind(params); err != nil {
		Response(c, 90400, err.Error())
		return
	}

	data, err := service.Certificate.GetDomainCertificateList(params.Name, params.Page, params.Limit)
	if err != nil {
		Response(c, 90500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": data,
	})
}

// DownloadDomainCertificate 下载证书
// @Summary 下载证书
// @Description 域名证书相关
// @Tags 域名证书相关
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path int true "域名ID"
// @Router /api/v1/certificate/{id} [get]
func (cert *certificate) DownloadDomainCertificate(c *gin.Context) {

	certId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		Response(c, 90500, err.Error())
		return
	}

	zip, name, err := service.Certificate.DownloadDomainCertificate(certId)
	if err != nil {
		Response(c, 90500, err.Error())
		return
	}

	zipFileName := fmt.Sprintf("%s.zip", name)
	c.Header("Content-Disposition", "attachment; filename="+zipFileName)

	c.Header("Content-Type", "application/zip")
	c.Data(http.StatusOK, "application/zip", zip.Bytes())
}
