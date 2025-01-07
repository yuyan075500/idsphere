package routers

import (
	"github.com/gin-gonic/gin"
	"ops-api/controller"
)

func initSettingsRouters(router *gin.Engine) {
	// 获取配置
	router.GET("/api/v1/settings", controller.Settings.GetSettings)

	settings := router.Group("/api/v1/settings")
	{
		// 上传 Logo
		settings.POST("/logoUpload", controller.Settings.UploadLogo)
		// 获取 Logo（不需要认证，不需要权限）
		settings.GET("/logo", controller.Settings.GetLogo)
		// 修改配置
		settings.PUT("", controller.Settings.UpdateSettings)
	}
}
