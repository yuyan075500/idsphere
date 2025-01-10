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
		// 发送邮箱测试
		settings.POST("/test/mailSend", controller.Settings.SendMail)
		// LDAP 登录测试
		settings.POST("/test/ldapLogin", controller.Settings.LdapLogin)
		// 发送短信测试
		settings.POST("/test/smsSend", controller.Settings.SendSms)
	}
}
