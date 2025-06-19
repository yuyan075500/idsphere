package routers

import (
	"github.com/gin-gonic/gin"
	"ops-api/controller"
)

// 初始化用户相关路由
func initUserRouters(router *gin.Engine) {
	// 获取用户列表（表格）
	router.GET("/api/v1/users", controller.User.GetUserList)

	user := router.Group("/api/v1/user")
	{
		// 新增用户
		user.POST("", controller.User.AddUser)
		// 修改用户
		user.PUT("", controller.User.UpdateUser)
		// 删除用户
		user.DELETE("/:id", controller.User.DeleteUser)
		// 重置用户密码（管理员）
		user.PUT("/reset_password", controller.User.UpdateUserPassword)
		// 重置用户MFA
		user.PUT("/reset_mfa/:id", controller.User.ResetUserMFA)
		// 获取用户列表（下拉框：分组用户管理）
		user.GET("/list", controller.User.GetUserListAll)
		// 用户头像上传
		user.POST("/avatarUpload", controller.User.UploadAvatar)
		// 从LDAP从步用户
		user.POST("/sync/ad", controller.User.UserSyncAd)
	}

	// 发送密码重置验证码
	router.POST("/api/v1/sms/reset_password", controller.User.GetVerificationCode)
	// 重置用户密码（用户自己）
	router.POST("/api/v1/reset_password", controller.User.UpdateSelfPassword)
}
