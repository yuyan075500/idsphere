package routers

import (
	"github.com/gin-gonic/gin"
	"ops-api/controller"
)

// 初始化账号相关路由
func initAccountRouters(router *gin.Engine) {
	// 获取账号列表（表格）
	router.GET("/api/v1/accounts", controller.Account.GetAccountList)

	// 批量新增账号
	router.POST("/api/v1/accounts", controller.Account.AddAccounts)

	account := router.Group("/api/v1/account")
	{
		// 新增账号
		account.POST("", controller.Account.AddAccount)
		// 删除删除
		account.DELETE("/:id", controller.Account.DeleteAccount)
		// 修改账号
		account.PUT("", controller.Account.UpdateAccount)
		// 批量修改账号所有者
		account.PUT("/owners", controller.Account.BatchUpdateAccountOwner)
		// 账号分享
		account.PUT("/users", controller.Account.UpdateAccountUser)
		// 获取账号密码
		account.GET("/password/:id", controller.Account.GetAccountPassword)
		// 更改密码
		account.PUT("/password", controller.Account.UpdatePassword)
		// 获取短信验证码
		account.POST("/code", controller.Account.GetSMSCode)
		// 获取短信验证码
		account.POST("/code_verification", controller.Account.CodeVerification)
	}
}
