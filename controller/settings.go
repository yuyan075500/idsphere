package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"ops-api/service"
	"path/filepath"
)

var Settings settings

type settings struct{}

// GetSettings 获取配置
// @Summary 获取配置
// @Description 配置相接口
// @Tags 配置相接口
// @Param Authorization header string true "Bearer 用户令牌"
// @Success 200 {string} json "{"code": 0, "data": []}"
// @Router /api/v1/settings [get]
func (s *settings) GetSettings(c *gin.Context) {

	data, err := service.Settings.GetAllSettingsWithParsedValues()
	if err != nil {
		Response(c, 90500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": data,
	})
}

// GetLogo 获取 Logo
// @Summary 获取 Logo
// @Description 配置相接口
// @Tags 配置相接口
// @Param Authorization header string true "Bearer 用户令牌"
// @Success 200 {string} json "{"code": 0, "path": logoPreview}"
// @Router /api/v1/settings/logo [get]
func (s *settings) GetLogo(c *gin.Context) {

	logoPreview, err := service.Settings.GetLogo()
	if err != nil {
		Response(c, 90500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"path": logoPreview,
	})
}

// UpdateSettings 修改配置
// @Summary 修改配置
// @Description 配置相接口
// @Tags 配置相接口
// @Param Authorization header string true "Bearer 用户令牌"
// @Param task body service.SettingsUpdate true "配置信息"
// @Success 200 {string} json "{"code": 0, "msg": "更新成功", "data": nil}"
// @Router /api/v1/settings [put]
func (s *settings) UpdateSettings(c *gin.Context) {
	var data = &service.SettingsUpdate{}

	if err := c.ShouldBind(&data); err != nil {
		Response(c, 90400, err.Error())
		return
	}

	// 更新
	result, err := service.Settings.UpdateSettingValues(data)
	if err != nil {
		Response(c, 90500, err.Error())
		return
	}

	CreateOrUpdateResponse(c, 0, "更新成功", result)
}

// UploadLogo 上传 Logo
// @Summary 上传 Logo
// @Description 配置相接口
// @Tags 配置相接口
// @Param Authorization header string true "Bearer 用户令牌"
// @Param logo formData file true "Logo"
// @Success 200 {string} json "{"code": 0, "path": logoPreview}"
// @Router /api/v1/settings/logoUpload [post]
func (s *settings) UploadLogo(c *gin.Context) {
	// 获取上传的Logo
	logo, err := c.FormFile("logo")
	if err != nil {
		Response(c, 90400, err.Error())
		return
	}

	// 拼接存储的路径（此路径为临时路径，在表单提交时会将图片移动到实际位置）
	logoPath := fmt.Sprintf("settings/logo/%s%s", uuid.New(), filepath.Ext(logo.Filename))

	// 执行上传
	logoPreview, err := service.Settings.UploadLogo(logoPath, logo)
	if err != nil {
		Response(c, 90500, err.Error())
		return
	}

	// 保存
	if err := service.Settings.UpdateSettingValue("logo", logoPath); err != nil {
		Response(c, 90500, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"path": logoPreview,
		"msg":  "Logo 上传成功",
	})
}