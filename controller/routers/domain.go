package routers

import (
	"github.com/gin-gonic/gin"
	"ops-api/controller"
)

func initDomainRouters(router *gin.Engine) {
	// 获取域名服务商列表
	router.GET("/api/v1/domain/providers", controller.Domain.GetDomainServiceProviderList)

	provider := router.Group("/api/v1/domain/provider")
	{
		// 新增域名服务商
		provider.POST("", controller.Domain.AddDomainServiceProvider)

		// 删除域名服务商
		provider.DELETE("/:id", controller.Domain.DeleteDomainServiceProvider)

		// 修改域名服务商
		provider.PUT("", controller.Domain.UpdateDomainServiceProvider)
	}
}
